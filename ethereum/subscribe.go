package ethereum

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (e *Ethereum) SubscribeBlockchain(address string) error {
	if e.ethWsClient == nil {
		return errors.New("ethereum websocket client is not initialized")
	}

	headers := make(chan *types.Header)
	subscribe, err := e.ethWsClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-subscribe.Err():
			return err
		case header := <-headers:
			block, err := e.ethClient.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				continue
			}

			for _, transaction := range block.Transactions() {
				if transaction == nil || transaction.To() == nil {
					continue
				}

				if transaction.To().String() == address {
					paidAmount, _ := e.WeiToEther(transaction.Value()).Float64()

					fmt.Printf("[Blockterminal] %v ETH paid to %v\n", paidAmount, transaction.To().String())
				}
			}
		}
	}
}
