package terminal

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/skatebord/blockterminal/wallets"
)

type DaisyTransferCommand struct {
	BasicCommand

	terminal *Terminal
}

func daisyTransferCommand(terminal *Terminal) *DaisyTransferCommand {
	return &DaisyTransferCommand{
		BasicCommand: newBasicCommand("daisytransfer", "Funnel funds through multiple wallets to create a daisy chain to a final destination", []string{
			"<from> (loaded wallet name)",
			"<recipient> (recipient address or loaded wallet name)",
			"<amount>",
			"<create_new_wallets_for_daisy_chain> (optional, will create a new wallet for each chain, default is false)",
			"<max_depth> (optional, will default to 10)",
			"<currency> (optional, will default to current chain default currency, e.g. USDT, USDC, DAI, etc.)",
			"<chain_wallets_path> (optional, will default to wallets/{random_name})",
			"<chain_wallet_password> (optional, will create a new password for the chain wallets)",
		}),
		terminal: terminal,
	}
}

func (c *DaisyTransferCommand) Execute(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("invalid number of arguments")
	}

	loadedWallet, err := c.terminal.wallets.GetLoadedWallet(args[0])
	if err != nil {
		return fmt.Errorf("failed to get loaded wallet: %v", err)
	}

	if loadedWallet.Chain() != c.terminal.chain.Name() {
		return fmt.Errorf("wallet is not on the current chain! Please select the correct chain")
	}

	amount := args[1]
	recipient := args[2]

	toAddress := ""
	if toWallet, err := c.terminal.wallets.GetLoadedWallet(recipient); err == nil {
		toAddress = toWallet.Address()
	} else {
		toAddress = recipient
	}

	createNewWallets := args[3]
	currency := args[4]
	maxDepth := args[5]
	_ = currency

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount")
	}

	depth := 10

	if maxDepth != "" {
		depth, err = strconv.Atoi(maxDepth)
		if err != nil {
			return fmt.Errorf("invalid max depth")
		}
	}

	createNewWalletsBool := false
	if createNewWallets == "true" {
		createNewWalletsBool = true
	}

	chainWalletPassword := args[6]
	if chainWalletPassword == "" {
		chainWalletPassword = "chain_wallet_password" // #TODO: make this a random password
	}

	chainWalletsPath := args[5]
	if chainWalletsPath == "" {
		chainWalletsPath = fmt.Sprintf("wallets/%s", uuid.New().String())
	}

	if !createNewWalletsBool {
		fmt.Println("Using currently loaded wallets for daisy chain")
	}

	fmt.Printf("Chain: %s\n", c.terminal.chain.Name())
	fmt.Printf("Depth: %d\n", depth)
	fmt.Printf("Chain wallets path: %s\n", chainWalletsPath)
	fmt.Printf("Chain wallet password: %s\n", chainWalletPassword)

	walletPrivateKey, err := loadedWallet.Unlock()
	if err != nil {
		return fmt.Errorf("failed to unlock wallet: %v", err)
	}
	_ = walletPrivateKey

	chainWallets := make([]wallets.Wallet, 0)

	if createNewWalletsBool {
		for i := 0; i < depth; i++ {
			newWallet := wallets.NewLoadedWallet(fmt.Sprintf("chain_%d", i), "", c.terminal.chain.Name())

			err = wallets.SaveWalletWithKeys(fmt.Sprintf("%s/chain_%d.json", chainWalletsPath, i), c.terminal.chain.Name(), newWallet,
				"", chainWalletPassword)
			if err != nil {
				return fmt.Errorf("failed to save chain wallet: %v", err)
			}

			chainWallets = append(chainWallets, newWallet)
		}
	} else {
		chainWallets = c.terminal.wallets.GetAllLoadedWalletsWithout(loadedWallet.Name())
	}

	if len(chainWallets) == 0 {
		return fmt.Errorf("no chain wallets found")
	}

	// start chain, send from wallet amount to first loaded wallet
	txHash, err := c.terminal.chain.SendAndConfirm(loadedWallet, chainWallets[0].Address(), amountFloat)
	if err != nil {
		return fmt.Errorf("failed to send and confirm: %v", err)
	}

	fmt.Printf("Chain #1: %s\n", txHash)

	for i := 1; i < depth; i++ {
		sendWallet := chainWallets[i]
		nextWallet := chainWallets[i+1]

		txHash, err := c.terminal.chain.SendAndConfirm(sendWallet, nextWallet.Address(), amountFloat)
		if err != nil {
			return fmt.Errorf("failed to send and confirm (last chain: %d, next chain: %d): %v", i, i+1, err)
		}

		fmt.Printf("Chain #%d: %s\n", i+1, txHash)
	}

	fmt.Println("Daisy chain! All wallets have been created and funded to the final chain wallet")

	loadedWallet, err = c.terminal.wallets.GetLoadedWallet(chainWallets[depth-1].Name())
	if err != nil {
		return fmt.Errorf("failed to get loaded wallet: %v", err)
	}

	// final chain wallet, send all to recipient
	txHash, err = c.terminal.chain.SendAndConfirm(loadedWallet, recipient, amountFloat)
	if err != nil {
		return fmt.Errorf("failed to send and confirm (final chain: %d, recipient: %s): %v", depth-1, toAddress, err)
	}

	fmt.Printf("Final Chain: %s\n", txHash)

	fmt.Println("Daisy chain complete! All funds have been sent to the recipient")

	return nil
}
