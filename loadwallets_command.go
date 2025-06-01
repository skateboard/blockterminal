package terminal

import (
	"fmt"

	"github.com/skatebord/blockterminal/wallets"
)

type LoadWalletsCommand struct {
	BasicCommand

	terminal *Terminal
}

func loadWalletsCommand(terminal *Terminal) *LoadWalletsCommand {
	return &LoadWalletsCommand{
		BasicCommand: newBasicCommand("loadwallets", "Load a directory of wallets from a path", []string{
			"<path>",
		}),
		terminal: terminal,
	}
}

func (c *LoadWalletsCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	loadedWallets, err := wallets.LoadWallets(args[0])
	if err != nil {
		return fmt.Errorf("failed to load wallet: %v", err)
	}

	// test for duplicate wallet names
	loadedWalletNames := make(map[string]bool)
	cleanedWallets := make([]*wallets.LoadedWallet, 0)
	for _, wallet := range loadedWallets {
		if loadedWalletNames[wallet.Name()] {
			return fmt.Errorf("wallet %s already loaded or a wallet by the same name exists", wallet.Name())
		}

		if wallet.Chain() != c.terminal.chain.Name() {
			fmt.Printf("Wallet %s is not on the current chain! Please select the correct chain\n", wallet.Name())
			continue
		}

		loadedWalletNames[wallet.Name()] = true
		cleanedWallets = append(cleanedWallets, wallet)
	}

	for _, wallet := range cleanedWallets {
		c.terminal.wallets.SaveLoadedWallet(wallet)
		c.terminal.chain.AddSubscribedAddress(wallet.Address())
	}

	fmt.Printf("Wallets loaded successfully\n")

	return nil
}
