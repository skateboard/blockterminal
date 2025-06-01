package terminal

import (
	"flag"
	"fmt"
	"os"

	"github.com/skatebord/blockterminal/wallets"
)

var (
	nodePath    = flag.String("node", "", "the node to automatically connect to")
	walletPath  = flag.String("wallet", "", "the wallet path to automatically load")
	walletsPath = flag.String("wallets", "", "the wallets path to automatically load a list of wallets")
)

func (t *Terminal) parseArgs() {
	if len(os.Args) == 0 {
		return
	}

	flag.Parse()

	if *nodePath != "" {
		err := t.connectNode(*nodePath)
		if err != nil {
			fmt.Println("Error connecting to node:", err)
		}
	}

	if *walletPath != "" {
		if t.chain == nil {
			fmt.Println("No chain selected, please use --node to connect to a node first")
			return
		}

		wallet, err := wallets.LoadWallet(*walletPath)
		if err != nil {
			fmt.Println("Error loading wallet:", err)
			return
		}

		t.wallets.SaveLoadedWallet(wallet)
		t.chain.AddSubscribedAddress(wallet.Address())

		fmt.Printf("Wallet %s loaded successfully\n", wallet.Name())
		fmt.Printf("%s\n", wallet.Address())
	}

	if *walletsPath != "" {
		if t.chain == nil {
			fmt.Println("No chain selected, please use --node to connect to a node first")
			return
		}

		wallets, err := wallets.LoadWallets(*walletsPath)
		if err != nil {
			fmt.Println("Error loading wallets:", err)
			return
		}

		for _, wallet := range wallets {
			t.wallets.SaveLoadedWallet(wallet)
			t.chain.AddSubscribedAddress(wallet.Address())
		}

		fmt.Printf("Wallets loaded successfully\n")
	}

}
