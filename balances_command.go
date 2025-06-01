package terminal

import (
	"fmt"
)

type BalancesCommand struct {
	BasicCommand

	terminal *Terminal
}

func balancesCommand(terminal *Terminal) *BalancesCommand {
	return &BalancesCommand{
		BasicCommand: newBasicCommand("balances", "Show balances of all loaded wallets", []string{}),
		terminal:     terminal,
	}
}

func (c *BalancesCommand) Execute(args []string) error {
	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	loadedWallets := c.terminal.wallets.GetAllLoadedWallets()

	for _, loadedWallet := range loadedWallets {
		if loadedWallet.Chain() != c.terminal.chain.Name() {
			fmt.Printf("Wallet %s is not on the current chain! Please select the correct chain\n", loadedWallet.Name())
			continue
		}

		balances, err := c.terminal.chain.GetBalance(loadedWallet.Address())
		if err != nil {
			fmt.Printf("Failed to get balance for %s: %v\n", loadedWallet.Name(), err)
			continue
		}

		fmt.Printf("Balances for %s:\n", loadedWallet.Name())
		fmt.Printf("Address: %s\n", loadedWallet.Address())
		fmt.Println("--------------------------------")
		for currency, balance := range balances {
			fmt.Printf("%s: %.2f\n", currency, balance)
		}
		fmt.Println("--------------------------------")
		fmt.Println()
	}

	return nil
}
