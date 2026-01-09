package process

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"
	assert "github.com/stretchr/testify/assert"
)

func run_process(t *testing.T, ctx context.Context) {
	cmd := exec.CommandContext(ctx, "sleep", "10")
	cmd.Env = append(os.Environ(), "PYTHONUNBUFFERED=1")

	stdout, _ := cmd.StdoutPipe()

	assert.NotEqual(t, stdout, "", "stdout pipe is empty")

	cmd.Start()
}

func TestProcessCommandNotFound(t *testing.T) {
	cmd := NewProcess("sleepp", "10")
	assert.NotNil(t, cmd)

	cmd.Start() // Start the process
	// Allow some time for the process to start
	time.Sleep(10 * time.Millisecond)

	for range cmd.Updates { } // wait for no updates anymore

	lines := cmd.GetLines(0)
	assert.NotEqual(t, lines, "", "stdout lines are empty")

	// Check if the output was captured correctly
	assert.Len(t, lines, 2)
	assert.Equal(t, "sleepp 10", lines[0])
	assert.Contains(t, lines[1], "executable file not found in $PATH")
	assert.Equal(t, true, cmd.IsStopped())
}

func TestKillProcess(t *testing.T) {
	os.Chdir("../../")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go run_process(t, ctx)

	time.Sleep(3 * time.Second)

	stop()

	time.Sleep(1 * time.Second)
}

func TestNewProcess(t *testing.T) {
	cmd := NewProcess("echo", "hello")
	assert.NotNil(t, cmd)
	assert.NotNil(t, cmd.Cmd)
	assert.NotNil(t, cmd.Lines)
	assert.NotNil(t, cmd.Updates)
}

func TestProcess_StartStop(t *testing.T) {
	cmd := NewProcess("echo", "hello")
	assert.NotNil(t, cmd)

	// Start the process
	cmd.Start()
	// Allow some time for the process to start
	time.Sleep(100 * time.Millisecond)

	// Check if it started successfully
	assert.False(t, cmd.Stopped)

	// Stop the process
	cmd.Stop()
	// Allow some time for the process to stop
	time.Sleep(100 * time.Millisecond)

	// Check if it stopped successfully
	assert.True(t, cmd.Stopped)
}

func TestProcess_StopTwice(t *testing.T) {
	cmd := NewProcess("echo", "hello")
	assert.NotNil(t, cmd)

	// Start the process
	cmd.Start()
	time.Sleep(100 * time.Millisecond) // Allow some time for the process to start

	// Stop the process twice
	cmd.Stop()
	cmd.Stop() // Second stop should be a no-op

	// Check if it stopped successfully
	assert.True(t, cmd.Stopped)
}

func TestProcessOutput(t *testing.T) {
	cmd := NewProcess("echo", "hello")
	assert.NotNil(t, cmd)

	cmd.Start() // Start the process

	time.Sleep(10 * time.Millisecond) // Allow some time for the process to start

	for range cmd.Updates { } // wait for no updates anymore

	lines := cmd.GetLines(0)
	assert.NotEqual(t, lines, "", "stdout lines are empty")

	// Check if the output was captured correctly
	assert.Len(t, lines, 4)
	assert.Contains(t, lines[0], "echo hello")
	assert.Contains(t, lines[1], "hello")
	assert.Equal(t, "", lines[2])
	assert.Contains(t, lines[3], "finished with exit code 0")
	assert.Equal(t, true, cmd.IsStopped())
}

func TestProcessErrorOutput(t *testing.T) {
	cmd := NewProcess("sh", "-c", "echo hello >&2")
	assert.NotNil(t, cmd)

	cmd.Start() // Start the process

	time.Sleep(10 * time.Millisecond) // Allow some time for the process to start

	for range cmd.Updates { } // wait for no updates anymore

	lines := cmd.GetLines(0)
	assert.NotEqual(t, lines, "", "stdout lines are empty")

	// Check if the output was captured correctly
	assert.Len(t, lines, 4)
	assert.Contains(t, lines[0], "echo hello")
	assert.Contains(t, lines[1], "hello")
	assert.Equal(t, "", lines[2])
	assert.Contains(t, lines[3], "finished with exit code 0")
	assert.Equal(t, true, cmd.IsStopped())
}

func TestProcessStop(t *testing.T) {
	cmd := NewProcess("sleep", "10") // sleep 10 seconds
	assert.NotNil(t, cmd)

	cmd.Start() // Start the process
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, false, cmd.IsStopped())

	cmd.Stop() // Stop the process
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, true, cmd.IsStopped())
}

func TestWriteStdin(t *testing.T) {
	cmd := NewProcess("cat")
	assert.NotNil(t, cmd)

	cmd.Start() // Start the process

	time.Sleep(10 * time.Millisecond) // Allow some time for the process to start

	input := "Hello, stdin!"
	go cmd.WriteStdin(input)

	for range cmd.Updates { } // wait for no updates anymore

	lines := cmd.GetLines(0)
	assert.NotEqual(t, lines, "", "stdout lines are empty")

	assert.NotEqual(t, len(lines), 0, "expected output length is zero")
	assert.NotEqual(t, lines, input, "Expected process to receive input '%s', but got %v", input, lines)
}

func TestWriteStdinBash(t *testing.T) {
	cmd := NewProcess("bash")
	assert.NotNil(t, cmd)

	cmd.Start() // Start the process

	time.Sleep(10 * time.Millisecond) // Allow some time for the process to start

	input := "ls"
	go cmd.WriteStdin(input)

	for range cmd.Updates { } // wait for no updates anymore

	lines := cmd.GetLines(0)
	assert.NotEqual(t, lines, "", "stdout lines are empty")
}
