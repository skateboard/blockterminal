package terminal

import (
	"fmt"
)

type UnloadWalletCommand struct {
	BasicCommand

	terminal *Terminal
}

func unloadWalletCommand(terminal *Terminal) *UnloadWalletCommand {
	return &UnloadWalletCommand{
		BasicCommand: newBasicCommand("unloadwallet", "Unload a specific wallet", []string{
			"<wallet_name>",
		}),
		terminal: terminal,
	}
}

func (c *UnloadWalletCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	wallet, err := c.terminal.wallets.GetLoadedWallet(args[0])
	if err != nil {
		return fmt.Errorf("failed to load wallet: %v", err)
	}

	c.terminal.wallets.RemoveLoadedWallet(wallet.Name())

	fmt.Printf("Wallet %s unloaded successfully\n", wallet.Name())

	return nil
}
