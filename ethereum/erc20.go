package ethereum

import (
	"math"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/skatebord/blockterminal/ethereum/contracts"
)

var (
	USDT_ADDRESS = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	USDC_ADDRESS = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	DAI_ADDRESS  = "0x6b175474e89094c44da98b954eedeac495271d0f"
)

type erc20Contract struct {
	address  common.Address
	contract *contracts.Erc20
}

func newErc20Contract(backend bind.ContractBackend, address string) (*erc20Contract, error) {
	contract, err := contracts.NewErc20(common.HexToAddress(address), backend)
	if err != nil {
		return nil, err
	}

	return &erc20Contract{
		address:  common.HexToAddress(address),
		contract: contract,
	}, nil
}

func (c *erc20Contract) Name() (string, error) {
	return c.contract.Name(&bind.CallOpts{})
}

func (c *erc20Contract) Decimals() (uint8, error) {
	decimals, err := c.contract.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}

	return uint8(decimals.Int64()), nil
}

func (c *erc20Contract) Balance(address string) (float64, error) {
	balance, err := c.contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		return 0, err
	}

	decimals, err := c.Decimals()
	if err != nil {
		return 0, err
	}

	return float64(balance.Uint64()) / math.Pow(10, float64(decimals)), nil
}

func (c *erc20Contract) Symbol() (string, error) {
	return c.contract.Symbol(&bind.CallOpts{})
}

func (c *erc20Contract) Address() string {
	return c.address.Hex()
}
