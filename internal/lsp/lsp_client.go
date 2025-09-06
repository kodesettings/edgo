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

	message2chan          map[int]chan string
	completionMessages    chan string
	definitionMessages    chan string
	referencesMessages    chan string
	signatureHelpMessages chan string
	hoverMessages         chan string
	otherMessages         chan string
	DiagnosticsChannel    chan string

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
	l.message2chan = make(map[int]chan string)
	l.completionMessages = make(chan string)
	l.referencesMessages = make(chan string)
	l.definitionMessages = make(chan string)
	l.signatureHelpMessages = make(chan string)
	l.hoverMessages = make(chan string)
	l.otherMessages = make(chan string)
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

func (this *LspClient) receiveDiagnostics() string {

	const LEN_HEADER = "Content-Length: "
	var messageSize int
	var responseMustBeNext bool
	var line string
	var err error
repeat:
	if messageSize != 0 && responseMustBeNext {
		buf := make([]byte, messageSize)
		_, err = io.ReadFull(this.reader, buf)
		if err != nil { slog.Error(err.Error()); goto repeat; }
		line = string(buf)
		messageSize = 0

		responseJSON := make(map[string]interface{})
		err = json.Unmarshal(buf, &responseJSON)
		if err != nil { slog.Error(err.Error()); goto repeat; }

		method, methodFound := responseJSON["method"]
		if methodFound && method.(string) == "textDocument/publishDiagnostics" {
			var dr DiagnosticResponse
			err = json.Unmarshal(buf, &dr)
			if err != nil { slog.Error(err.Error()); goto repeat; }
			return line
		}

		if value, idFound := responseJSON["id"]; idFound {
			if _, ok := value.(float64); ok {
				return line
			}
		}
	} else {
		line, err = this.reader.ReadString('\n') // it stuck sometimes
		if err != nil { slog.Error("[445 lsp]", "err", err.Error()); goto repeat; }
	}

	line = strings.TrimSuffix(line, "\r\n")

	if strings.HasPrefix(line, LEN_HEADER) {
		sizeStr := strings.TrimPrefix(line, LEN_HEADER)
		msize, _ := strconv.Atoi(sizeStr)
		messageSize = msize
		responseMustBeNext = false
		goto repeat;
	}

	if line == "" {
		responseMustBeNext = true
		goto repeat;
	}

	return ""
}

func (l *LspClient) receiveLoop() {
repeat:
	message := l.receiveDiagnostics()
	slog.Info("recv:", "json", message)

	if strings.Contains(message,"publishDiagnostics") {
		var dr DiagnosticResponse
		err := json.Unmarshal([]byte(message), &dr)
		if err != nil { slog.Error(err.Error()); goto repeat; }
		l.file2diagnostic[dr.Params.Uri] = dr.Params
		l.DiagnosticsChannel <- message
		goto repeat;
	}
	if strings.Contains(message,"workspace/applyEdit") {
		l.otherMessages <- message
		goto repeat;
	}

	responseJSON := make(map[string]interface{})
	err := json.Unmarshal([]byte(message), &responseJSON)
	if err != nil { slog.Error(err.Error()); goto repeat; }

	if value, found := responseJSON["id"]; found { // json has id
		if id, ok := value.(float64); ok {
			channel, foundRequest := l.message2chan[int(id)]
			if foundRequest {
				channel <- message
			} else  {
				//skip message
			}
		}
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

	l.message2chan[id] = l.otherMessages
	l.send(initializeRequest)

	response, err := WaitForRequest[interface{}](l.otherMessages, 3000)

	delete(l.message2chan, id)

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
	id := this.id

	request := BaseRequest{
		ID: id, JSONRPC: "2.0", Method:  "textDocument/hover",
		Params: Params {
			TextDocument: TextDocument { URI: "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.message2chan[id] = this.hoverMessages
	this.send(request)

	response, err := WaitForRequest[HoverResponse](this.hoverMessages, 1000)

	delete(this.message2chan, id)
	return response, err
}

func (this *LspClient) Completion(file string, line int, character int) (CompletionResponse, error) {
	this.id++
	id := this.id

	request := BaseRequest{
		ID: id, JSONRPC: "2.0", Method:  "textDocument/completion",
		Params: Params{
			TextDocument: TextDocument { URI:  "file://" + file },
			Position: Position { Line: line, Character: character },
			Context: Context { TriggerKind: 1 },
		},
	}

	this.message2chan[id] = this.completionMessages
	this.send(request)

	response, err := WaitForRequest[CompletionResponse](this.completionMessages, 1000)

	delete(this.message2chan, id)
	return response, err
}

func (this *LspClient) Definition(file string, line int, character int) (DefinitionResponse, error) {
	this.id++
	id := this.id

	request := DefinitionRequest{
		ID: this.id, JSONRPC: "2.0", Method:  "textDocument/definition",
		Params: DefinitionParams {
			TextDocument: TextDocument{ URI: "file://" + file },
			Position: Position{ Line: line, Character: character },
		},
	}

	this.message2chan[id] = this.definitionMessages
	this.send(request)

	response, err := WaitForRequest[DefinitionResponse](this.definitionMessages, 1000)

	delete(this.message2chan, id)
	return response, err
}

func (this *LspClient) SignatureHelp(file string, line int, character int) (SignatureHelpResponse, error) {
	this.id++
	id := this.id

	request := BaseRequest{
		ID: id, JSONRPC: "2.0", Method:  "textDocument/signatureHelp",
		Params: Params {
			TextDocument: TextDocument { URI: "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.message2chan[id] = this.signatureHelpMessages
	this.send(request)

	response, err := WaitForRequest[SignatureHelpResponse](this.signatureHelpMessages, 1000)

	delete(this.message2chan, id)
	return response, err
}

func (this *LspClient) References(file string, line int, character int) (ReferencesResponse, error) {
	this.id++
	id := this.id

	request := BaseRequest{
		ID: id, JSONRPC: "2.0", Method:  "textDocument/references",
		Params: Params{
			TextDocument: TextDocument{ URI: "file://" + file },
			Position: Position{ Line: line, Character: character },
			Context: Context{ IncludeDeclaration: false },
		},
	}

	this.message2chan[id] = this.referencesMessages
	this.send(request)

	response, err := WaitForRequest[ReferencesResponse](this.referencesMessages, 3000)

	delete(this.message2chan, id)
	return response, err
}


func (this *LspClient) PrepareRename(file string, line int, character int) (PrepareRenameResponse, error) {
	this.id++
	id := this.id

	request := PrepareRenameRequest {
		ID: id, JSONRPC: "2.0", Method:  "textDocument/prepareRename",
		Params: Params{
			TextDocument: TextDocument { URI:  "file://" + file },
			Position: Position { Line: line, Character: character },
		},
	}

	this.message2chan[id] = this.otherMessages
	this.send(request)

	response, err := WaitForRequest[PrepareRenameResponse](this.otherMessages, 10000)

	delete(this.message2chan, id)
	return response, err
}

func (this *LspClient) Rename(file string, newname string, line int, character int) (RenameResponse, error) {
	this.id++
	id := this.id

	request := RenameRequest{
		ID: id,  JSONRPC: "2.0", Method:  "textDocument/rename",
		Params: RenameParams {
			NewName: newname,
			Position: Position { Line: line, Character: character },
			TextDocument: TextDocument { URI:  "file://" + file },
		},
	}

	this.message2chan[id] = this.otherMessages
	this.send(request)

	response, err := WaitForRequest[RenameResponse](this.otherMessages, 10000)

	delete(this.message2chan, id)
	return response, err
}

func (this *LspClient) CodeAction(file string, spc int, spl int, epc int, epl int) (CodeActionResponse, error) {
	this.id++
	id := this.id

	request := CodeActionRequest {
		ID: id,  JSONRPC: "2.0", Method: "textDocument/codeAction",
		Params: CodeActionParams {
			TextDocument: TextDocument { URI:  "file://" + file },
			Context: Context{ Only: []string{"refactor"}, TriggerKind: 1 },
			Range: RequestRange{
				Start: Position{ Line: spl, Character: spc},
				End: Position{ Line: epl, Character: epc},
			},
		},
	}

	this.message2chan[id] = this.otherMessages
	this.send(request)

	response, err := WaitForRequest[CodeActionResponse](this.otherMessages, 10000)

	delete(this.message2chan, id)
	return response, err
}

func (this *LspClient) Command(command Command) (CommandResponse, error) {
	this.id++
	id := this.id

	request := CommandRequest {
		ID: id, JSONRPC: "2.0", Method: "workspace/executeCommand",
		Params: command,
	}

	this.send(request)

	jsonData := <- this.otherMessages

	var response CommandResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil { slog.Error("Error parsing JSON:", "err", err.Error()) }
	return response, err
}

func (this *LspClient) ApplyEdit(key int) {
	request := ApplyEditRequest {
		ID: key,  JSONRPC: "2.0",
		Result:  Applied { true } ,
	}
	this.send(request)
}
