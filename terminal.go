package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skatebord/blockterminal/wallets"
)

var (
	VERSION = "0.0.1"
)

type Terminal struct {
	chain Chain

	commandRegistry *CommandRegistry

	currentWallet wallets.Wallet
}

func NewTerminal() *Terminal {
	terminal := &Terminal{
		commandRegistry: NewCommandRegistry(),
	}
	terminal.initialize()

	return terminal
}

func (t *Terminal) initialize() {
	t.commandRegistry.RegisterCommands([]Command{
		createNodeCommand(),
		connectNodeCommand(t),
		createWalletCommand(t),
		loadWalletCommand(t),
		balancesCommand(t),
		transferCommand(t),
		infoCommand(t),
		exitCommand(),
		helpCommand(t.commandRegistry),
	})
}

func (t *Terminal) Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Welcome to Blockterminal v%s\n", VERSION)
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
	walletName := ""
	chainName := ""

	if t.currentWallet != nil {
		walletName = t.currentWallet.Name()
	}

	if t.chain != nil {
		chainName = t.chain.GetNodeName()
	}

	if walletName == "" && chainName == "" {
		return "> "
	}

	return fmt.Sprintf("%s@%s> ", walletName, chainName)
}

func (t *Terminal) SetChain(chain Chain) {
	t.chain = chain
}

func (t *Terminal) SetWallet(wallet wallets.Wallet) {
	t.currentWallet = wallet
}
