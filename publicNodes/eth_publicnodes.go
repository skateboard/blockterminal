package publicnodes

var EthereumPublicNodes = map[string]PublicNode{
	"eth_mainnet": {
		Name:      "ethMainnet",
		RPC:       "https://ethereum-rpc.publicnode.com",
		WS:        "wss://ethereum-rpc.publicnode.com",
		ChainType: "ethereum",
	},
	"eth_sepolia": {
		Name:      "ethSepolia",
		RPC:       "https://ethereum-sepolia-rpc.publicnode.com",
		WS:        "wss://ethereum-sepolia-rpc.publicnode.com",
		ChainType: "ethereum",
	},
}
