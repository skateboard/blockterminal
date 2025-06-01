package terminal

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/skatebord/blockterminal/wallets"
)

type Chain interface {
	Name() string
	GetNodeName() string

	Send(fromWallet wallets.Wallet, toAddress string, amount float64) (string, error)
	SendAndConfirm(fromWallet wallets.Wallet, toAddress string, amount float64) (string, error)

	AddSubscribedAddress(address string)
	SubscribeBlockchain()
	GetBalance(address string) (map[string]float64, error)
}

type ChainConfig struct {
	ChainType string `json:"chain"`
	Name      string `json:"name"`
	Rpc       string `json:"rpc"`
	Ws        string `json:"ws"`
}

func LoadChainConfigs() (map[string]*ChainConfig, error) {
	files, err := os.ReadDir("nodes")
	if err != nil {
		return nil, err
	}

	configs := make(map[string]*ChainConfig)
	for _, file := range files {
		config, err := LoadChainConfig(fmt.Sprintf("nodes/%s", file.Name()))
		if err != nil {
			return nil, err
		}

		configs[config.Name] = config
	}

	return configs, nil
}

func SaveChainConfig(path string, config *ChainConfig) error {
	json, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, json, 0644)
}

func LoadChainConfig(path string) (*ChainConfig, error) {
	j, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ChainConfig
	err = json.Unmarshal(j, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
