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
			"<to>",
			"<amount>",
			"<currency> (optional, will default to current chain default currency)",
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

	if c.terminal.currentWallet == nil {
		return fmt.Errorf("no wallet selected, please select a wallet first")
	}

	to := args[0]
	amount := args[1]
	// currency := args[2]

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount")
	}

	txHash, err := c.terminal.chain.Send(c.terminal.currentWallet, to, amountFloat)
	if err != nil {
		return err
	}

	fmt.Printf("Transferred %s to %s\n", amount, to)
	fmt.Printf("%s\n", txHash)

	return nil
}
