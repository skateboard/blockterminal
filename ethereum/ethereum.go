package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type Ethereum struct {
	nodeName    string
	ethClient   *ethclient.Client
	ethWsClient *ethclient.Client

	erc20Contracts map[string]*erc20Contract
}

func New(name string, rpc string, wsRpc string) (*Ethereum, error) {
	ethClient, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	ethWsClient, err := ethclient.Dial(wsRpc)
	if err != nil {
		return nil, err
	}

	erc20Contracts := make(map[string]*erc20Contract)
	erc20Contracts[USDT_ADDRESS], err = newErc20Contract(ethClient, USDT_ADDRESS)
	if err != nil {
		return nil, err
	}

	erc20Contracts[USDC_ADDRESS], err = newErc20Contract(ethClient, USDC_ADDRESS)
	if err != nil {
		return nil, err
	}

	erc20Contracts[DAI_ADDRESS], err = newErc20Contract(ethClient, DAI_ADDRESS)
	if err != nil {
		return nil, err
	}

	return &Ethereum{nodeName: name, ethClient: ethClient, ethWsClient: ethWsClient, erc20Contracts: erc20Contracts}, nil
}

func (e *Ethereum) Name() string {
	return "ethereum"
}

func (e *Ethereum) GetNodeName() string {
	return e.nodeName
}

func (e *Ethereum) GetBalance(address string) (map[string]float64, error) {
	balance, err := e.ethClient.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}
	f, _ := e.WeiToEther(balance).Float64()

	balances := map[string]float64{
		"Ethereum (ETH)": f,
	}

	for _, erc20Contract := range e.erc20Contracts {
		balance, err := erc20Contract.Balance(address)
		if err != nil {
			continue
		}

		name, err := erc20Contract.Name()
		if err != nil {
			continue
		}

		symbol, err := erc20Contract.Symbol()
		if err != nil {
			continue
		}

		balances[fmt.Sprintf("%s (%s)", name, symbol)] = balance
	}

	return balances, nil
}

func (e *Ethereum) EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func (e *Ethereum) WeiToEther(wei *big.Int) *big.Float {
	weiFloat := new(big.Float).SetInt(wei)
	eth := new(big.Float).Quo(weiFloat, big.NewFloat(params.Ether))
	return eth
}
