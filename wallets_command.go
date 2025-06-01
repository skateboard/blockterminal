package terminal

import (
	"fmt"
)

type WalletsCommand struct {
	BasicCommand

	terminal *Terminal
}

func walletsCommand(terminal *Terminal) *WalletsCommand {
	return &WalletsCommand{
		BasicCommand: newBasicCommand("wallets", "Show all loaded wallets", []string{}),
		terminal:     terminal,
	}
}

func (c *WalletsCommand) Execute(args []string) error {
	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	loadedWallets := c.terminal.wallets.GetAllLoadedWallets()

	fmt.Printf("Loaded wallets (%d):\n", len(loadedWallets))
	for _, loadedWallet := range loadedWallets {
		fmt.Printf("Wallet: %s\n", loadedWallet.Name())
		fmt.Printf("Address: %s\n", loadedWallet.Address())
		fmt.Printf("Chain: %s\n", loadedWallet.Chain())
		fmt.Println("--------------------------------")
	}

	return nil
}
