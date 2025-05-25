package terminal

import (
	"fmt"

	"github.com/skatebord/blockterminal/ethereum"
	"github.com/skatebord/blockterminal/wallets"
)

type CreateWalletCommand struct {
	BasicCommand

	terminal *Terminal
}

func createWalletCommand(terminal *Terminal) *CreateWalletCommand {
	return &CreateWalletCommand{
		BasicCommand: newBasicCommand("createwallet", "Create a new wallet", []string{
			"<path>",
			"<wallet_name>",
			"<password>",
			"<chain> (optional, if not provided, the terminal chain will be used)",
		}),
		terminal: terminal,
	}
}

func (c *CreateWalletCommand) Execute(args []string) error {
	if len(args) != 3 && len(args) != 4 {
		return fmt.Errorf("invalid number of arguments")
	}

	path := args[0]
	walletName := args[1]
	password := args[2]

	chainName := ""
	if len(args) == 4 {
		chainName = args[3]
	} else {
		if c.terminal.chain == nil {
			return fmt.Errorf("no chain selected, please provide a chain")
		}
		chainName = c.terminal.chain.Name()
	}
	_ = chainName

	var wallet wallets.Wallet
	var privateKey string
	switch chainName {
	case "ethereum":
		ethWallet, err := ethereum.NewWallet()
		if err != nil {
			return fmt.Errorf("failed to create wallet: %v", err)
		}

		wallet = wallets.NewLoadedWallet(walletName, ethWallet.Address())

		privateKey = ethWallet.SaveToHex()
	}

	err := wallets.SaveWalletWithKeys(path, chainName, wallet, privateKey, password)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %v", err)
	}

	fmt.Printf("Wallet %s created successfully\n", walletName)
	fmt.Printf("%s\n", wallet.Address())

	return nil
}
