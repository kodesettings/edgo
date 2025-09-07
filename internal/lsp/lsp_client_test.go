package lsp

import (
	"os"
	"path"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestLspClientFindProcess(t *testing.T) {
	lsp := LspClient{}
	started := lsp.Start("gopls")
	defer lsp.stop()

	assert.Equal(t, true, started, "lsp not started")
	assert.NotNil(t, lsp.Cmd, "cmd is nil")

	pid := lsp.Cmd.Process.Pid
	assert.NotEqual(t, pid, 0, "lsp pid is zero")

	process, err := os.FindProcess(pid)
	assert.Nil(t, err, "error finding Cmd with id %d: %s\n", process.Pid, err)
}

func TestLspClientInitialize(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")
	defer lsp.stop()

	pid := lsp.Cmd.Process.Pid
	assert.NotEqual(t, pid, 0, "lsp pid is zero")

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	assert.Equal(t, true, lsp.IsReady, "expected lsp to be ready, got false")
}

func TestLspClientHover(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")
	defer lsp.stop()

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.Hover(file, 50-1, 8)
	assert.Nil(t, err, "error occured")

	expected := "func (this *LspClient) DidOpen(file string, text *string)"
	got := response.Result.Contents.Value
	assert.Equal(t, expected, got, "expected lsp hover result to be %s, got something else %s", expected, got)
}

func TestLspClientCompletion(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")
	defer lsp.stop()

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.Completion(file, 71-1, 8)
	assert.Nil(t, err, "error occured")

	expected := "DidChange"
	assert.Equal(t, 4, len(response.Result.Items), "items must be 4")
	assert.Equal(t, expected, response.Result.Items[0].Label, "expected lsp completion result to be %s", expected)
}

func TestLspClientDefinition(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")
	defer lsp.stop()

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.Definition(file, 92-1, 8)
	assert.Nil(t, err, "error occured")

	assert.Equal(t, 1, len(response.Result), "items must be 1")
	assert.NotEqual(t, "", response.Result[0].URI, "expected lsp definition uri result to be non empty")
}

func TestLspClientSignatureHelp(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")
	defer lsp.stop()

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.SignatureHelp(file, 112-1, 8)
	assert.Nil(t, err, "error occured")

	expected := "DidOpen(file string, text *string)"
	assert.Equal(t, 1, len(response.Result.Signatures), "items must be 1")
	assert.Equal(t, expected, response.Result.Signatures[0].Label, "expected lsp signature help result to be %s", expected)
}

func TestLspClientReferences(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")
	defer lsp.stop()

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.References(file, 133-1, 8)
	assert.Nil(t, err, "error occured")

	assert.Equal(t, 6, len(response.Result), "items must be 6")
	assert.NotEqual(t, "", response.Result[0].URI, "expected lsp references result uri to be non empty")
}
