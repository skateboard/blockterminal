package terminal

import "fmt"

type BalancesCommand struct {
	BasicCommand

	terminal *Terminal
}

func balancesCommand(terminal *Terminal) *BalancesCommand {
	return &BalancesCommand{
		BasicCommand: newBasicCommand("balances", "Show balances of the current wallet", []string{}),
		terminal:     terminal,
	}
}

func (c *BalancesCommand) Execute(args []string) error {
	if c.terminal.currentWallet == nil {
		return fmt.Errorf("no wallet selected, please select a wallet first")
	}

	balances, err := c.terminal.chain.GetBalance(c.terminal.currentWallet.Address())
	if err != nil {
		return err
	}

	fmt.Println("Balances:")
	for currency, balance := range balances {
		fmt.Printf("%s: %.2f\n", currency, balance)
	}

	return nil
}
