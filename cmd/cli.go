package main

import (
	terminal "github.com/skatebord/blockterminal"
)

func main() {
	// rpc := flag.String("rpc", "https://ethereum-sepolia-rpc.publicnode.com", "the rpc url")
	// wsRpc := flag.String("wsRpc", "wss://ethereum-sepolia-rpc.publicnode.com", "the ws rpc url")
	// chain := flag.String("chain", "ethereum", "the chain")

	// flag.Parse()

	// var terminalChain terminal.Chain
	// var err error

	// switch *chain {
	// case "ethereum":
	// 	terminalChain, err = ethereum.NewEthereum(*rpc, *wsRpc)
	// }

	// if err != nil {
	// 	log.Fatalf("Failed to create terminal chain: %v", err)
	// 	return
	// }

	terminal.NewTerminal().Run()
}
