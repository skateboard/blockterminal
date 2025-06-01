package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	btHttp "github.com/skatebord/blockterminal/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

type Ethereum struct {
	nodeName    string
	ethClient   *ethclient.Client
	ethWsClient *ethclient.Client

	subscribedAddresses map[string]bool

	erc20Contracts map[string]*erc20Contract
}

func New(name string, rpcUrl string, wsURL string, http *btHttp.Http) (*Ethereum, error) {
	rpcClient, err := rpc.DialOptions(context.Background(), rpcUrl, rpc.WithHTTPClient(http.GetHttpClient()))
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)

	var ethWsClient *ethclient.Client
	if wsURL != "" {
		wsDialer := websocket.Dialer{
			NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return http.GetDialer().Dial(network, addr)
			},
			HandshakeTimeout: 45 * time.Second,
		}

		wsRPCClient, err := rpc.DialOptions(context.Background(), wsURL, rpc.WithWebsocketDialer(wsDialer))
		if err != nil {
			return nil, err
		}
		ethWsClient = ethclient.NewClient(wsRPCClient)
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

	return &Ethereum{nodeName: name,
		ethClient:           ethClient,
		ethWsClient:         ethWsClient,
		erc20Contracts:      erc20Contracts,
		subscribedAddresses: make(map[string]bool),
	}, nil
}

func (e *Ethereum) AddSubscribedAddress(address string) {
	if e.subscribedAddresses[address] {
		return
	}

	e.subscribedAddresses[address] = true
	return
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

		if balance == 0 {
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

func (e *Ethereum) isErc20Contract(address string) bool {
	_, ok := e.erc20Contracts[address]
	return ok
}

func (e *Ethereum) GetErc20Contract(address string) (*erc20Contract, error) {
	contract, ok := e.erc20Contracts[address]
	if !ok {
		return nil, fmt.Errorf("contract not found")
	}
	return contract, nil
}
