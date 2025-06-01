package terminal

import (
	"fmt"

	"github.com/skatebord/blockterminal/ethereum"
	publicnodes "github.com/skatebord/blockterminal/publicNodes"
)

type ConnectNodeCommand struct {
	BasicCommand

	terminal *Terminal
}

func connectNodeCommand(terminal *Terminal) *ConnectNodeCommand {
	return &ConnectNodeCommand{
		BasicCommand: newBasicCommand("connectnode", "Connect to a node", []string{
			"<node_name>",
		}),
		terminal: terminal,
	}
}

func (c *ConnectNodeCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	nodeName := args[0]

	var chainConfig *ChainConfig
	if publicnodes.PublicNodes[nodeName] {
		publicNodeConfig := publicnodes.EthereumPublicNodes[nodeName]

		chainConfig = &ChainConfig{
			Name:      publicNodeConfig.Name,
			ChainType: publicNodeConfig.ChainType,
			Rpc:       publicNodeConfig.RPC,
			Ws:        publicNodeConfig.WS,
		}
	} else {
		chainConfigs, err := LoadChainConfigs()
		if err != nil {
			return fmt.Errorf("failed to load chain configs: %v", err)
		}

		c, ok := chainConfigs[nodeName]
		if !ok {
			return fmt.Errorf("node %s not found", nodeName)
		}

		chainConfig = c
	}

	var chain Chain
	var err2 error
	switch chainConfig.ChainType {
	case "ethereum":
		chain, err2 = ethereum.New(chainConfig.Name, chainConfig.Rpc, chainConfig.Ws)
	}

	if err2 != nil {
		return fmt.Errorf("failed to create chain: %v", err2)
	}

	c.terminal.SetChain(chain)

	fmt.Printf("Connected to node %s successfully\n", nodeName)

	return nil
}
