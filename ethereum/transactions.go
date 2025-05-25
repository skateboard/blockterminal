package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/skatebord/blockterminal/wallets"
)

func (e *Ethereum) Send(fromWallet wallets.Wallet, toAddress string, amount float64) (string, error) {
	privateKeyHex, err := fromWallet.Unlock()
	if err != nil {
		return "", err
	}

	unlockedFromWallet, err := LoadWalletFromHex(privateKeyHex)
	if err != nil {
		return "", err
	}

	fromAddress := unlockedFromWallet.Address()

	nonce, err := e.ethClient.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		return "", err
	}

	gas := big.NewInt(21000) // in units
	gasPrice, err := e.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	cost := gas.Mul(gas, gasPrice)

	sendAmount := e.EtherToWei(big.NewFloat(amount))
	sendAmount = sendAmount.Sub(sendAmount, cost)

	var data []byte
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), sendAmount, uint64(21000), gasPrice, data)

	chainID, err := e.ethClient.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), unlockedFromWallet.PrivateKey)
	if err != nil {
		return "", err
	}

	err = e.ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
