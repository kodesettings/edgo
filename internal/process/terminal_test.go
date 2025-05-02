package process

import (
	"fmt"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestTerminalStart(t *testing.T) {
	terminal, err := NewTerminal()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer terminal.Pty.Close()

	// Example: Execute "ls" command on the terminal
	err = terminal.ExecuteCommand("ls")
	assert.Nil(t, err, "error received in terminal start: %s", err)

	// Access all output lines
	fmt.Println("Output Lines:")
	terminal.mutex.Lock()
	for _, line := range terminal.Lines {
		fmt.Print(line)
	}
	terminal.mutex.Unlock()
}
