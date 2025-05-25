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

	if c.terminal.currentWallet == nil {
		return fmt.Errorf("no wallet selected, please select a wallet first")
	}

	fmt.Printf("Chain: %s\n", c.terminal.chain.Name())
	fmt.Printf("Wallet: %s\n", c.terminal.currentWallet.Name())
	fmt.Printf("Address: %s\n", c.terminal.currentWallet.Address())
	return nil
}
