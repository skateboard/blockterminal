package terminal

import (
	"fmt"
)

type UnloadWalletsCommand struct {
	BasicCommand

	terminal *Terminal
}

func unloadWalletsCommand(terminal *Terminal) *UnloadWalletsCommand {
	return &UnloadWalletsCommand{
		BasicCommand: newBasicCommand("unloadwallets", "Unload all wallets", []string{}),
		terminal:     terminal,
	}
}

func (c *UnloadWalletsCommand) Execute(args []string) error {
	c.terminal.wallets.Clear()

	fmt.Printf("All wallets unloaded successfully\n")

	return nil
}
