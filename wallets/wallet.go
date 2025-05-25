package wallets

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type Wallet interface {
	Name() string
	Chain() string
	Address() string
	Unlock() (string, error)
}

type LoadedWallet struct {
	chain   string
	name    string
	address string

	encryptedKey string
	salt         string
}

func NewLoadedWallet(name, address, chain string) *LoadedWallet {
	return &LoadedWallet{
		chain:   chain,
		name:    name,
		address: address,
	}
}

func (w *LoadedWallet) Chain() string {
	return w.chain
}

func (w *LoadedWallet) Name() string {
	return w.name
}

func (w *LoadedWallet) Address() string {
	return w.address
}

func (w *LoadedWallet) Unlock() (string, error) {
	fmt.Printf("Enter password: ")

	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	password := strings.TrimSpace(string(bytePassword))

	return decryptPrivateKey(w.encryptedKey, w.salt, password)
}

func LoadWallet(path string) (*LoadedWallet, error) {
	j, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wallet map[string]any
	err = json.Unmarshal(j, &wallet)
	if err != nil {
		return nil, err
	}

	encryptedKey, ok := wallet["encrypted_privatekey"].(string)
	if !ok {
		return nil, errors.New("encrypted_privatekey is not a string")
	}

	salt, ok := wallet["salt"].(string)
	if !ok {
		return nil, errors.New("salt is not a string")
	}

	return &LoadedWallet{
		chain:        wallet["chain"].(string),
		name:         wallet["name"].(string),
		address:      wallet["address"].(string),
		encryptedKey: encryptedKey,
		salt:         salt,
	}, nil
}

func SaveWalletWithKeys(path, chain string, wallet Wallet, privateKey string, password string) error {
	encryptedKey, salt, err := encryptPrivateKey(privateKey, password)
	if err != nil {
		return err
	}

	j := map[string]interface{}{
		"chain":                chain,
		"name":                 wallet.Name(),
		"address":              wallet.Address(),
		"encrypted_privatekey": encryptedKey,
		"salt":                 salt,
	}

	json, err := json.Marshal(j)
	if err != nil {
		return err
	}

	return os.WriteFile(path, json, 0644)
}
