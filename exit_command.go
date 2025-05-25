package terminal

import (
	"fmt"
	"os"
)

type ExitCommand struct {
	BasicCommand
}

func exitCommand() *ExitCommand {
	return &ExitCommand{
		BasicCommand: newBasicCommand("exit", "Exit the terminal", []string{}),
	}
}

func (c *ExitCommand) Execute(args []string) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}
