package publicnodes

// This package is for a list of public nodes that can be used to connect to a blockchain.
// We currently have:
// 	Ethereum Mainnet
// 	Ethereum Sepolia

var (
	// made for easy indexing
	PublicNodes = map[string]bool{
		"eth_mainnet": true,
		"eth_sepolia": true,
	}
)

type PublicNode struct {
	Name      string
	ChainType string
	RPC       string
	WS        string
}
