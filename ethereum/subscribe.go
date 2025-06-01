package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (e *Ethereum) SubscribeBlockchain() {
	if e.ethWsClient == nil {
		return
	}

	headers := make(chan *types.Header)
	subscribe, err := e.ethWsClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return
	}

	for {
		select {
		case err := <-subscribe.Err():
			fmt.Println(err)
			return
		case header := <-headers:
			block, err := e.ethClient.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				continue
			}

			for _, transaction := range block.Transactions() {
				if transaction == nil || transaction.To() == nil {
					continue
				}

				address := transaction.To().String()

				if e.isErc20Contract(address) {
					erc20Contract, err := e.GetErc20Contract(address)
					if err != nil {
						continue
					}

					receipt, err := e.ethClient.TransactionReceipt(context.Background(), transaction.Hash())
					if err != nil {
						continue
					}

					log := receipt.Logs[0]

					transfer, err := erc20Contract.DecodeTransferEvent(log)
					if err != nil {
						continue
					}

					if !e.subscribedAddresses[transfer.To.String()] {
						continue
					}

					symbol, _ := erc20Contract.Symbol()
					paidAmount, _ := e.WeiToEther(transaction.Value()).Float64()

					fmt.Printf("[Blockterminal] %v %v paid to %v\n", paidAmount, symbol, transfer.To.String())
					continue
				}

				if !e.subscribedAddresses[address] {
					continue
				}

				paidAmount, _ := e.WeiToEther(transaction.Value()).Float64()

				fmt.Printf("[Blockterminal] %v ETH paid to %v\n", paidAmount, transaction.To().String())
			}
		}
	}
}
