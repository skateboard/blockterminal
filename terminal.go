package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skatebord/blockterminal/ethereum"
	"github.com/skatebord/blockterminal/wallets"
)

var (
	VERSION = "0.0.1"
)

type Terminal struct {
	chain Chain

	commandRegistry *CommandRegistry

	wallets *wallets.Wallets
}

func NewTerminal() *Terminal {
	terminal := &Terminal{
		commandRegistry: NewCommandRegistry(),
		wallets:         wallets.NewWallets(),
	}
	terminal.initialize()
	terminal.parseArgs()

	return terminal
}

func (t *Terminal) initialize() {
	t.commandRegistry.RegisterCommands([]Command{
		createNodeCommand(),
		connectNodeCommand(t),
		nodesCommand(t),
		createWalletCommand(t),
		loadWalletCommand(t),
		loadWalletsCommand(t),
		balancesCommand(t),
		balanceCommand(t),
		walletsCommand(t),
		transferCommand(t),
		infoCommand(t),
		exitCommand(),
		helpCommand(t.commandRegistry),
	})
}

func (t *Terminal) Run() {
	fmt.Printf("Welcome to Blockterminal v%s\n", VERSION)

	reader := bufio.NewReader(os.Stdin)
	for {
		// Print the custom prompt
		fmt.Printf("%s", t.buildPrompt())

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)

		commandName := strings.Split(input, " ")[0]
		commandArgs := strings.Split(input, " ")[1:]

		err = t.commandRegistry.RunCommand(commandName, commandArgs)
		if err != nil {
			fmt.Println("Error running command:", err)
		}
	}
}

func (t *Terminal) buildPrompt() string {
	chainName := ""

	if t.chain != nil {
		chainName = t.chain.GetNodeName()
	}

	if chainName == "" {
		return "> "
	}

	return fmt.Sprintf("%s> ", chainName)
}

func (t *Terminal) SetChain(chain Chain) {
	t.chain = chain
}

func (t *Terminal) connectNode(nodePath string) error {
	chainConfig, err := LoadChainConfig(nodePath)
	if err != nil {
		return fmt.Errorf("error loading chain config: %v", err)
	}

	var chain Chain

	switch chainConfig.ChainType {
	case "ethereum":
		chain, err = ethereum.New(chainConfig.Name, chainConfig.Rpc, chainConfig.Ws)
		if err != nil {
			return fmt.Errorf("error creating chain: %v", err)
		}
	}

	t.SetChain(chain)

	fmt.Printf("Connected to %s\n", chainConfig.Name)

	return nil

}
