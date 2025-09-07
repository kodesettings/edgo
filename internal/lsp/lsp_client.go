package lsp

import (
	"bufio"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"bytes"
	"time"
	"log/slog"
)

type LspClient struct {
	Lang   string
	Cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stop   context.CancelFunc
	reader *bufio.Reader

	IsReady bool // flag for lsp initialization

	userMessages       chan string
	DiagnosticsChannel chan string

	id              int
	file2diagnostic map[string]DiagnosticParams
}

func (l *LspClient) Start(cmd string, args ...string) bool {
	slog.Info("starting lsp", cmd, strings.Join(args," "))

	_, err := exec.LookPath(cmd)
	if err != nil { slog.Info("lsp not found ", "cmd", cmd); return false }

	ctx, stop := signal.NotifyContext(context.Background(), os.Kill)
	l.Cmd = exec.CommandContext(ctx, cmd, args...)
	l.stop = stop

	stdin, err := l.Cmd.StdinPipe()
	if err != nil { slog.Info(err.Error()); return false }
	l.stdin = stdin

	stdout, err := l.Cmd.StdoutPipe()
	if err != nil { slog.Info(err.Error()); return false }
	l.stdout = stdout

	// starting lsp Cmd async
	startError := l.Cmd.Start()
	if startError != nil {
		slog.Error("error starting lsp ", "err", startError.Error())
		return false
	} else {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		slog.Info("lsp started success", cmd, strings.Join(args," "))
	}

	l.reader = bufio.NewReader(stdout)
	l.userMessages = make(chan string, 10)
	l.DiagnosticsChannel = make(chan string, 10)
	l.file2diagnostic = make(map[string]DiagnosticParams)
	l.IsReady = true

	go l.receiveLoop()

	return true
}

func (this *LspClient) send(o interface{})  {
	m, err := json.Marshal(o)
	if err != nil { panic(err) }

	slog.Info("send:", "json", string(m))
	message := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(m), m)
	_, err = this.stdin.Write([]byte(message))
	if err != nil { slog.Error(err.Error()) }
}

func (this *LspClient) receiveDiagnostics() (DiagnosticResponse, string) {
	var dr DiagnosticResponse = DiagnosticResponse{}

	content_len, _, err := this.reader.ReadLine()
	if err != nil { slog.Error("io readline", "error", err.Error()); return dr, "" }

	nbytes_str := bytes.TrimPrefix(content_len, []byte("Content-Length: "))
	nbytes_conv, err := strconv.Atoi(string(nbytes_str))
	if err != nil { slog.Error("strconv", "error", err.Error()); return dr, "" }

	_, _, err = this.reader.ReadLine() // reading line separator
	buf := make([]byte, int(nbytes_conv)) // allocate parsed content length

	_, err = io.ReadFull(this.reader, buf)
	if err != nil { slog.Error("io read", "error", err.Error()); return dr, "" }

	// wrapping the extracted json body into a stream and
	// trying to decode the value into a struct
	reader := bufio.NewReader(bytes.NewReader(buf))
	err = json.NewDecoder(reader).Decode(&dr)
	if err != nil { slog.Error("json", "error", err.Error()); return dr, "" }

	// returning both the diagnostic response and the string representation
	return dr, string(buf)
}

func (l *LspClient) receiveLoop() {
repeat:
	dr, message := l.receiveDiagnostics()
	slog.Info("recv:", "json", message)

	if strings.Contains(message,"publishDiagnostics") {
		l.file2diagnostic[dr.Params.Uri] = dr.Params
		l.DiagnosticsChannel <- message
	} else if strings.Contains(message, "result") {
		l.userMessages <- message
	}

	goto repeat;
}

func (this *LspClient) GetDiagnostic(filename string) (DiagnosticParams, bool) {
	d, found := this.file2diagnostic[filename]
	return  d, found
}

func WaitForRequest[T any](channel chan string, timeout int) (T, error) {
	var response T
	var err error

	select {
	case jsonData := <- channel:
		err = json.Unmarshal([]byte(jsonData), &response)
		if err != nil { slog.Error("Error parsing JSON:", "err", err.Error()) }

	case <-time.After(time.Duration(timeout) * time.Millisecond):
		err = fmt.Errorf("Timeout")
	}

	return response, err
}

func (l *LspClient) Init(dir string) {
	l.id = 0
	id := l.id

	initializeRequest := InitializeRequest{
		ID: id, JSONRPC: "2.0",
		Method: "initialize",
		Params: InitializeParams{
			RootURI: "file://" + dir, RootPath: dir,
			WorkspaceFolders: []WorkspaceFolder{{ Name: "edgo", URI:  "file://" + dir }},
			Capabilities: capabilities,
			ClientInfo: ClientInfo{ Name: "edgo",Version: "1.0.0"},
		},
	}

	l.send(initializeRequest)
	response, err := WaitForRequest[interface{}](l.userMessages, 3000)

	if response == "" || err != nil {
		slog.Info("cant get initialize response from lsp server")
		l.IsReady = false
		return
	}

	initializedRequest := InitializedRequest{
		JSONRPC: "2.0", Method:  "initialized", Params:  struct{}{},
	}
	l.send(initializedRequest)

	slog.Info("lsp initialized")
	l.IsReady = true
}

func (this *LspClient) DidOpen(file string, text *string) {
	didOpenRequest := DidOpenRequest{
		JSONRPC: "2.0",  Method:  "textDocument/didOpen",
		Params: DidOpenTextDocumentParams{
			TextDocument: TextDocument{
				LanguageID: this.Lang,
				Text:       *text,
				URI:        "file://" + file,
				Version:    1,
			},
		},
	}

	this.send(didOpenRequest)
}

func (this *LspClient) DidChange(file string, text *string, version int) {
	didChangeRequest := DidChangeRequest{
		JSONRPC: "2.0",  Method:  "textDocument/didChange",
		Params: DidChangeTextDocumentParams{
			ContentChanges: []TextDocumentContentChangeEvent{{
//				Range: Range{
//					Start: Position{ Line: spl, Character: spc},
//					End: Position{ Line: epl, Character: epc},
//				},
				Text: *text,
			}},
			TextDocument: VersionedTextDocumentIdentifier{
				URI:        "file://" + file,
				Version:    version,
			},
		},
	}

	this.send(didChangeRequest)
}

func (this *LspClient) DidClose(file string) {
	didCloseRequest := DidCloseRequest{
		JSONRPC: "2.0",  Method:  "textDocument/didClose",
		Params: DidCloseTextDocumentParams{
			TextDocument: TextDocumentIdentifier{
				URI:        "file://" + file,
			},
		},
	}

	this.send(didCloseRequest)
}

func (this *LspClient) Hover(file string, line int, character int) (HoverResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/hover",
		Params: Params {
			TextDocument: TextDocument { URI: "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[HoverResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) Completion(file string, line int, character int) (CompletionResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/completion",
		Params: Params{
			TextDocument: TextDocument { URI:  "file://" + file },
			Position: Position { Line: line, Character: character },
			Context: Context { TriggerKind: 1 },
		},
	}

	this.send(request)
	response, err := WaitForRequest[CompletionResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) Definition(file string, line int, character int) (DefinitionResponse, error) {
	this.id++

	request := DefinitionRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/definition",
		Params: DefinitionParams {
			TextDocument: TextDocument{ URI: "file://" + file },
			Position: Position{ Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[DefinitionResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) SignatureHelp(file string, line int, character int) (SignatureHelpResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/signatureHelp",
		Params: Params {
			TextDocument: TextDocument { URI: "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[SignatureHelpResponse](this.userMessages, 1000)

	return response, err
}

func (this *LspClient) References(file string, line int, character int) (ReferencesResponse, error) {
	this.id++

	request := BaseRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/references",
		Params: Params{
			TextDocument: TextDocument{ URI: "file://" + file },
			Position: Position{ Line: line, Character: character },
			Context: Context{ IncludeDeclaration: false },
		},
	}

	this.send(request)
	response, err := WaitForRequest[ReferencesResponse](this.userMessages, 3000)

	return response, err
}


func (this *LspClient) PrepareRename(file string, line int, character int) (PrepareRenameResponse, error) {
	this.id++

	request := PrepareRenameRequest {
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/prepareRename",
		Params: Params{
			TextDocument: TextDocument { URI:  "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.send(request)
	response, err := WaitForRequest[PrepareRenameResponse](this.userMessages, 10000)

	return response, err
}

func (this *LspClient) Rename(file string, newname string, line int, character int) (RenameResponse, error) {
	this.id++

	request := RenameRequest{
		ID: this.id,  JSONRPC: "2.0", Method:  "textDocument/rename",
		Params: RenameParams {
			NewName: newname,
			Position: Position { Line: line, Character: character },
			TextDocument: TextDocument { URI:  "file://" + file },
		},
	}

	this.send(request)
	response, err := WaitForRequest[RenameResponse](this.userMessages, 10000)

	return response, err
}

func (this *LspClient) CodeAction(file string, spc int, spl int, epc int, epl int) (CodeActionResponse, error) {
	this.id++

	request := CodeActionRequest {
		ID: this.id,  JSONRPC: "2.0", Method: "textDocument/codeAction",
		Params: CodeActionParams {
			TextDocument: TextDocument { URI:  "file://" + file },
			Context: Context{ Only: []string{"refactor"}, TriggerKind: 1 },
			Range: RequestRange{
				Start: Position{ Line: spl, Character: spc},
				End: Position{ Line: epl, Character: epc},
			},
		},
	}

	this.send(request)
	response, err := WaitForRequest[CodeActionResponse](this.userMessages, 10000)

	return response, err
}
