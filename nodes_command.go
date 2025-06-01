package terminal

import (
	"fmt"
)

type NodesCommand struct {
	BasicCommand

	terminal *Terminal
}

func nodesCommand(terminal *Terminal) *NodesCommand {
	return &NodesCommand{
		BasicCommand: newBasicCommand("nodes", "Show all saved nodes", []string{}),
		terminal:     terminal,
	}
}

func (c *NodesCommand) Execute(args []string) error {
	chainConfigs, err := LoadChainConfigs()
	if err != nil {
		return fmt.Errorf("failed to load chain configs: %v", err)
	}

	fmt.Printf("Saved nodes (%d):\n", len(chainConfigs))
	for _, chainConfig := range chainConfigs {
		fmt.Printf("Node: %s\n", chainConfig.Name)
		fmt.Printf("RPC: %s\n", chainConfig.Rpc)
		fmt.Printf("WS: %s\n", chainConfig.Ws)
		fmt.Printf("Chain: %s\n", chainConfig.ChainType)
		fmt.Println("--------------------------------")
	}

	return nil
}
