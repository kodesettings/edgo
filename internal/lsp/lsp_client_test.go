package lsp

import (
	. "github.com/vipmax/edgo/internal/logger"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestLspClientStart(t *testing.T) {
	lsp := LspClient{}
	started := lsp.Start("gopls")
	defer lsp.Stop()

	assert.Equal(t, started, true, "lsp not started")
	assert.NotNil(t, lsp.Cmd, "cmd is nil")

	pid := lsp.Cmd.Process.Pid
	fmt.Println("lsp pid is", pid)

	process, err := os.FindProcess(pid)
	assert.Nil(t, err, "error finding Cmd with id %d: %s\n", process.Pid, err)
	assert.NotEqual(t, lsp.isStopped, true, "expected lsp not to be stopped")
}

func TestLspClientStop(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")

	pid := lsp.Cmd.Process.Pid
	fmt.Println("lsp pid is", pid)

	lsp.Stop()

	assert.Equal(t, lsp.isStopped, true, "expected lsp to be stopped, got false")
}

func TestLspClientInitialize(t *testing.T) {
	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")

	pid := lsp.Cmd.Process.Pid
	fmt.Println("lsp pid is", pid)

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	assert.Equal(t, lsp.IsReady, true, "expected lsp to be ready, got false")
}

func TestLspClientHover(t *testing.T) {
	err := os.Setenv("EDGO_LOG", "edgo.log")
	Log.Start()

	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "internal","lsp", "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.Hover(file, 75-1, 7)
	assert.Nil(t, err, "error occured")

	fmt.Println(response, err)

	expected := "func (*Logger).Start()"
	got := response.Result.Contents.Value
	assert.Equal(t, got, expected, "expected lsp hover result to be %s, got something else %s", expected, got)
}

func TestLspClientCompletion(t *testing.T) {
	err := os.Setenv("EDGO_LOG", "edgo.log")
	Log.Start()

	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.Completion(file, 100-1, 8)
	assert.Nil(t, err, "error occured")

	fmt.Println(response, err)

	expected := "Start"
	assert.Equal(t, len(response.Result.Items), 1, "items must be 1")
	assert.Equal(t, response.Result.Items[0].Label, expected, "expeccted lsp completion result to be %s", expected)
}

func TestLspClientDefinition(t *testing.T) {
	err := os.Setenv("EDGO_LOG", "edgo.log")
	Log.Start()

	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.Definition(file, 124-1, 8)
	assert.Nil(t, err, "error occured")

	fmt.Println(response, err)

	expectedSuffix := "logger/logger.go"
	containsSuffix := strings.Contains(response.Result[0].URI, expectedSuffix)

	assert.Equal(t, len(response.Result), 1, "items must be 1")
	assert.Equal(t, containsSuffix, false, "expeccted lsp definition result to be %s", expectedSuffix)
}

func TestLspClientSignatureHelp(t *testing.T) {
	err := os.Setenv("EDGO_LOG", "edgo.log")
	Log.Start()

	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.SignatureHelp(file, 156-1, 21)
	assert.Nil(t, err, "error occured")

	fmt.Println(response, err)

	expected := "Join"
	assert.Equal(t, len(response.Result.Signatures), 1, "items must be 1")
	assert.Equal(t, response.Result.Signatures[0].Label, expected, "expeccted lsp definition result to be %s", expected)
}

func TestLspClientReferences(t *testing.T) {
	err := os.Setenv("EDGO_LOG", "edgo.log")
	Log.Start()

	lsp := LspClient{Lang: "go"}
	lsp.Start("gopls")

	currentDir, _ := os.Getwd()
	lsp.Init(currentDir)

	file := path.Join(currentDir, "internal","lsp", "lsp_client_test.go")
	text, _ := os.ReadFile(file)
	stringtext := string(text)
	lsp.DidOpen(file, &stringtext)

	response, err := lsp.References(file, 175-1, 2)
	assert.Nil(t, err, "error occured")

	fmt.Println(response, err)
	// todo fix, something wrong with base dir
}
