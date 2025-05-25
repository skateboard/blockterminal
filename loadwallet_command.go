package terminal

import (
	"fmt"

	"github.com/skatebord/blockterminal/wallets"
)

type LoadWalletCommand struct {
	BasicCommand

	terminal *Terminal
}

func loadWalletCommand(terminal *Terminal) *LoadWalletCommand {
	return &LoadWalletCommand{
		BasicCommand: newBasicCommand("loadwallet", "Load a wallet", []string{
			"<path>",
		}),
		terminal: terminal,
	}
}

func (c *LoadWalletCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	wallet, err := wallets.LoadWallet(args[0])
	if err != nil {
		return fmt.Errorf("failed to load wallet: %v", err)
	}

	if wallet.Chain() != c.terminal.chain.Name() {
		return fmt.Errorf("wallet is not on the current chain! Please select the correct chain first")
	}

	_, err = wallet.Unlock()
	if err != nil {
		return fmt.Errorf("failed to unlock wallet: %v", err)
	}

	c.terminal.SetWallet(wallet)

	fmt.Printf("Wallet %s loaded successfully\n", wallet.Name())
	fmt.Printf("%s\n", wallet.Address())

	go func() {
		err := c.terminal.chain.SubscribeBlockchain(wallet.Address())
		if err != nil {
			fmt.Printf("failed to subscribe to blockchain: %v\n", err)
			return
		}
	}()

	return nil
}
