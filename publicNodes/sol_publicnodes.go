package publicnodes

var SolanaPublicNodes = map[string]PublicNode{
	"sol_mainnet": {
		Name:      "solMainnet",
		RPC:       "https://solana-rpc.publicnode.com",
		WS:        "wss://solana-rpc.publicnode.com",
		ChainType: "solana",
	},
	"sol_testnet": {
		Name:      "solTestnet",
		RPC:       "https://solana-testnet-rpc.publicnode.com",
		WS:        "wss://solana-testnet-rpc.publicnode.com",
		ChainType: "solana",
	},
}
