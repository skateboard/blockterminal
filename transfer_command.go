package terminal

import (
	"fmt"
	"strconv"
)

type TransferCommand struct {
	BasicCommand

	terminal *Terminal
}

func transferCommand(terminal *Terminal) *TransferCommand {
	return &TransferCommand{
		BasicCommand: newBasicCommand("transfer", "Transfer funds", []string{
			"<from> (loaded wallet name)",
			"<to> (loaded wallet name or address)",
			"<amount>",
			"<currency> (optional, will default to current chain default currency, e.g. USDT, USDC, DAI, etc.)",
		}),
		terminal: terminal,
	}
}

func (c *TransferCommand) Execute(args []string) error {
	if len(args) != 2 && len(args) != 3 {
		return fmt.Errorf("invalid number of arguments")
	}

	if c.terminal.chain == nil {
		return fmt.Errorf("no chain selected, please select a chain first")
	}

	loadedWallet, err := c.terminal.wallets.GetLoadedWallet(args[0])
	if err != nil {
		return fmt.Errorf("failed to get loaded wallet: %v", err)
	}

	if loadedWallet.Chain() != c.terminal.chain.Name() {
		return fmt.Errorf("wallet is not on the current chain! Please select the correct chain")
	}

	to := args[1]

	toAddress := ""
	if toWallet, err := c.terminal.wallets.GetLoadedWallet(to); err == nil {
		toAddress = toWallet.Address()
	} else {
		toAddress = to
	}

	amount := args[2]
	// currency := args[2]

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount")
	}

	txHash, err := c.terminal.chain.Send(loadedWallet, toAddress, amountFloat)
	if err != nil {
		return err
	}

	fmt.Printf("Transferred %s to %s\n", amount, toAddress)
	fmt.Printf("%s\n", txHash)

	return nil
}
