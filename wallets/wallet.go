package wallets

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

type Wallet interface {
	Name() string
	Chain() string
	Address() string
	Unlock() (string, error)
}

type Wallets struct {
	loadedWallets map[string]Wallet
}

func NewWallets() *Wallets {
	return &Wallets{
		loadedWallets: make(map[string]Wallet),
	}
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

func LoadWallets(path string) ([]*LoadedWallet, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("path is not a directory")
	}

	loadedWallets := make([]*LoadedWallet, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		wallet, err := LoadWallet(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		loadedWallets = append(loadedWallets, wallet)
	}

	return loadedWallets, nil
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

	j := map[string]any{
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

func (w *Wallets) SaveLoadedWallet(loadedWallet *LoadedWallet) {
	w.loadedWallets[loadedWallet.name] = loadedWallet
}

func (w *Wallets) GetLoadedWallet(name string) (Wallet, error) {
	loadedWallet, ok := w.loadedWallets[name]
	if !ok {
		return nil, errors.New("wallet not found")
	}

	return loadedWallet, nil
}

func (w *Wallets) RemoveLoadedWallet(name string) {
	delete(w.loadedWallets, name)
}

func (w *Wallets) GetAllLoadedWallets() []Wallet {
	loadedWallets := make([]Wallet, 0)
	for _, loadedWallet := range w.loadedWallets {
		loadedWallets = append(loadedWallets, loadedWallet)
	}

	return loadedWallets
}

func (w *Wallets) GetAllLoadedWalletsWithout(name string) []Wallet {
	loadedWallets := make([]Wallet, 0)
	for _, loadedWallet := range w.loadedWallets {
		if loadedWallet.Name() == name {
			continue
		}

		loadedWallets = append(loadedWallets, loadedWallet)
	}

	return loadedWallets
}

func (w *Wallets) Len() int {
	return len(w.loadedWallets)
}

func (w *Wallets) Clear() {
	w.loadedWallets = make(map[string]Wallet)
}
