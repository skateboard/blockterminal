package terminal

import (
	"fmt"
	"os"
	"path/filepath"

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
			"<path> (must be a directory)",
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

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get path info: %v", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	walletPath := filepath.Join(path, walletName)

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

	var privateKey string
	var address string
	switch chainName {
	case "ethereum":
		ethWallet, err := ethereum.NewWallet()
		if err != nil {
			return fmt.Errorf("failed to create wallet: %v", err)
		}

		address = ethWallet.Address()
		privateKey = ethWallet.SaveToHex()
	}

	if address == "" || privateKey == "" {
		return fmt.Errorf("failed to create wallet, please try again")
	}

	wallet := wallets.NewLoadedWallet(walletName, address, chainName)

	err = wallets.SaveWalletWithKeys(walletPath, chainName, wallet, privateKey, password)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %v", err)
	}

	fmt.Printf("Wallet %s created successfully\n", walletName)
	fmt.Printf("%s\n", wallet.Address())

	return nil
}
