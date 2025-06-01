package terminal

import "fmt"

type InfoCommand struct {
	BasicCommand

	terminal *Terminal
}

func infoCommand(terminal *Terminal) *InfoCommand {
	return &InfoCommand{
		BasicCommand: newBasicCommand("info", "Show information about the current environment", []string{}),
		terminal:     terminal,
	}
}

func (c *InfoCommand) Execute(args []string) error {
	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	fmt.Printf("Chain: %s\n", c.terminal.chain.Name())
	fmt.Printf("Loaded Wallets: %d\n", c.terminal.wallets.Len())
	return nil
}
