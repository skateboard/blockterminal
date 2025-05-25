package terminal

import (
	"fmt"
)

type CreateNodeCommand struct {
	BasicCommand
}

func createNodeCommand() *CreateNodeCommand {
	return &CreateNodeCommand{
		BasicCommand: newBasicCommand("createnode", "Create a new node", []string{
			"<chain>",
			"<node_name>",
			"<rpc_url>",
			"<ws_url> (optional)",
		}),
	}
}

func (c *CreateNodeCommand) Execute(args []string) error {
	if len(args) != 3 && len(args) != 4 {
		return fmt.Errorf("invalid number of arguments")
	}

	chain := args[0]
	nodeName := args[1]
	rpcUrl := args[2]

	wsUrl := ""
	if len(args) == 4 {
		wsUrl = args[3]
	}

	err := SaveChainConfig(fmt.Sprintf("nodes/%s.json", nodeName), &ChainConfig{
		ChainType: chain,
		Name:      nodeName,
		Rpc:       rpcUrl,
		Ws:        wsUrl,
	})
	if err != nil {
		return fmt.Errorf("failed to save chain config: %v", err)
	}

	fmt.Printf("Node %s created successfully\n", nodeName)
	return nil
}
