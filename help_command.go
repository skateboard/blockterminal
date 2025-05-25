package terminal

import (
	"fmt"
	"strings"
)

type HelpCommand struct {
	BasicCommand

	commandRegistry *CommandRegistry
}

func helpCommand(commandRegistry *CommandRegistry) *HelpCommand {
	return &HelpCommand{
		BasicCommand:    newBasicCommand("help", "Show help", []string{}),
		commandRegistry: commandRegistry,
	}
}

func (c *HelpCommand) Execute(args []string) error {
	fmt.Println("Available commands:")
	for _, command := range c.commandRegistry.GetCommands() {
		fmt.Printf("%s %s - %s\n", command.Name(), strings.Join(command.Args(), " "), command.Description())
	}
	return nil
}
