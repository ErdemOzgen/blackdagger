package executor

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
)

// TimerExecutor wraps a CommandExecutor to add execution timing
type TimerExecutor struct {
	commandExecutor *CommandExecutor
}

func (e *TimerExecutor) Run() error {
	startTime := time.Now()
	err := e.commandExecutor.Run()
	executionTime := time.Since(startTime)
	// Here you can print the execution time or handle it as needed
	fmt.Printf("Execution Time: %v\n", executionTime)
	//log.L.Warnf("Execution Time: %v\n", executionTime)
	log.Printf("Execution Time \"%v\"\n", executionTime)

	return err
}

func (e *TimerExecutor) SetStdout(out io.Writer) {
	e.commandExecutor.SetStdout(out)
}

func (e *TimerExecutor) SetStderr(out io.Writer) {
	e.commandExecutor.SetStderr(out)
}

func (e *TimerExecutor) Kill(sig os.Signal) error {
	return e.commandExecutor.Kill(sig)
}

func CreateTimerExecutor(ctx context.Context, step *dag.Step) (Executor, error) {
	// Create the underlying command executor
	commandExecutor, err := CreateCommandExecutor(ctx, step)
	if err != nil {
		return nil, err
	}

	return &TimerExecutor{
		commandExecutor: commandExecutor.(*CommandExecutor),
	}, nil
}

func init() {
	//Register("", CreateCommandExecutor) // Keep existing registrations
	//Register("command", CreateCommandExecutor)
	// Register the timer executor
	Register("timer", CreateTimerExecutor)
}
