// terminal/ethereum/wallets.go
package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
}

// NewWallet generates a new Ethereum wallet
func NewWallet() (*Wallet, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return &Wallet{PrivateKey: priv}, nil
}

// SaveToHex returns the private key as a hex string
func (w *Wallet) SaveToHex() string {
	return hex.EncodeToString(crypto.FromECDSA(w.PrivateKey))
}

// LoadWalletFromHex loads a wallet from a hex-encoded private key
func LoadWalletFromHex(hexkey string) (*Wallet, error) {
	bytes, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}
	priv, err := crypto.ToECDSA(bytes)
	if err != nil {
		return nil, err
	}
	return &Wallet{PrivateKey: priv}, nil
}

// Address returns the Ethereum address for the wallet
func (w *Wallet) Address() string {
	return crypto.PubkeyToAddress(w.PrivateKey.PublicKey).Hex()
}
