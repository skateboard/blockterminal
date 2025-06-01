package terminal

import (
	"fmt"
)

type BalanceCommand struct {
	BasicCommand

	terminal *Terminal
}

func balanceCommand(terminal *Terminal) *BalanceCommand {
	return &BalanceCommand{
		BasicCommand: newBasicCommand("balance", "Show balance of a loaded wallet", []string{
			"<wallet_name>",
		}),
		terminal: terminal,
	}
}

func (c *BalanceCommand) Execute(args []string) error {
	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	walletName := args[0]

	loadedWallet, err := c.terminal.wallets.GetLoadedWallet(walletName)
	if err != nil {
		return fmt.Errorf("failed to load wallet: %v", err)
	}

	if loadedWallet.Chain() != c.terminal.chain.Name() {
		return fmt.Errorf("wallet is not on the current chain! Please select the correct chain first")
	}

	balances, err := c.terminal.chain.GetBalance(loadedWallet.Address())
	if err != nil {
		return err
	}

	fmt.Printf("Balances for %s:\n", loadedWallet.Name())
	fmt.Printf("Address: %s\n", loadedWallet.Address())
	fmt.Println("--------------------------------")
	for currency, balance := range balances {
		fmt.Printf("%s: %.2f\n", currency, balance)
	}
	fmt.Println("--------------------------------")
	fmt.Println()

	return nil
}
