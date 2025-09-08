package lsp

import (
	"os"
	"path"
	"testing"
	suite "github.com/stretchr/testify/suite"
)

type LspTestSuite struct {
	suite.Suite
	lsp LspClient
	started bool
}

// in order for 'go test' to run this suite, we need to create a
// normal test case that runs the suite
func TestLspTestSuite(t *testing.T) {
	suite.Run(t, new(LspTestSuite))
}

// make sure that lsp client is initialized for the entire suite
func (l *LspTestSuite) SetupSuite() {
	l.lsp = LspClient{Lang: "go"}
	l.started = l.lsp.Start("gopls")
}

func (l *LspTestSuite) TestLspClientFindProcess() {
	l.Equal(true, l.started, "lsp not started")
	l.NotNil(l.lsp.Cmd, "cmd is nil")

	pid := l.lsp.Cmd.Process.Pid
	l.NotEqual(pid, 0, "lsp pid is zero")

	process, err := os.FindProcess(pid)
	l.Nil(err, "error finding Cmd with id %d: %s\n", process.Pid, err)
}

//func (l *LspTestSuite) TestLspClientInitialize() {
//	pid := l.lsp.Cmd.Process.Pid
//	l.NotEqual(pid, 0, "lsp pid is zero")
//
//	currentDir, _ := os.Getwd()
//	l.lsp.Init(currentDir)
//
//	l.Equal(true, l.lsp.IsReady, "expected lsp to be ready, got false")
//}

func (l *LspTestSuite) TestLspClientHover() {
	currentDir, _ := os.Getwd()
	l.lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	l.lsp.DidOpen(file, &stringtext)

	response, err := l.lsp.Hover(file, 56-1, 8)
	l.Nil(err, "error occured")

	expected := "func (l *LspClient) DidOpen(file string, text *string)"
	got := response.Result.Contents.Value
	l.Equal(expected, got, "expected lsp hover result to be %s, got something else %s", expected, got)
}

func (l *LspTestSuite) TestLspClientCompletion() {
	currentDir, _ := os.Getwd()
	l. lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	l.lsp.DidOpen(file, &stringtext)

	response, err := l.lsp.Completion(file, 73-1, 8)
	l.Nil(err, "error occured")

	expected := "DiagnosticsChannel"
	l.Equal(9, len(response.Result.Items), "items must be 9")
	l.Equal(expected, response.Result.Items[0].Label, "expected lsp completion result to be %s", expected)
}

func (l *LspTestSuite) TestLspClientDefinition() {
	currentDir, _ := os.Getwd()
	l.lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	l.lsp.DidOpen(file, &stringtext)

	response, err := l.lsp.Definition(file, 90-1, 8)
	l.Nil(err, "error occured")

	l.Equal(1, len(response.Result), "items must be 1")
	l.NotEqual("", response.Result[0].URI, "expected lsp definition uri result to be non empty")
}

func (l *LspTestSuite) TestLspClientSignatureHelp() {
	currentDir, _ := os.Getwd()
	l.lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	l.lsp.DidOpen(file, &stringtext)

	response, err := l.lsp.SignatureHelp(file, 106-1, 8)
	l.Nil(err, "error occured")

	expected := "DidOpen(file string, text *string)"
	l.Equal(1, len(response.Result.Signatures), "items must be 1")
	l.Equal(expected, response.Result.Signatures[0].Label, "expected lsp signature help result to be %s", expected)
}

func (l *LspTestSuite) TestLspClientReferences() {
	currentDir, _ := os.Getwd()
	l.lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	l.lsp.DidOpen(file, &stringtext)

	response, err := l.lsp.References(file, 123-1, 8)
	l.Nil(err, "error occured")

	l.Equal(6, len(response.Result), "items must be 6")
	l.NotEqual("", response.Result[0].URI, "expected lsp references result uri to be non empty")
}
