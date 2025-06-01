package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

type Http struct {
	dialer     proxy.Dialer
	httpClient *http.Client
}

func New(useTor bool, port int) (*Http, error) {
	if useTor {
		torClient, torDialer, err := newTorHTTPClient(port)
		if err != nil {
			return nil, err
		}

		return &Http{dialer: torDialer, httpClient: torClient}, nil
	}

	httpClient, err := newHttpClient()
	if err != nil {
		return nil, err
	}

	// default dialer
	dialer := &net.Dialer{
		Timeout: 30 * time.Second,
	}

	return &Http{dialer: dialer, httpClient: httpClient}, nil
}

func (h *Http) GetDialer() proxy.Dialer {
	return h.dialer
}

func (h *Http) GetHttpClient() *http.Client {
	return h.httpClient
}

func newHttpClient() (*http.Client, error) {
	return &http.Client{
		Timeout: 30 * time.Second,
	}, nil
}

func newTorHTTPClient(port int) (*http.Client, proxy.Dialer, error) {
	torDialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), nil, proxy.Direct)
	if err != nil {
		return nil, nil, err
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return torDialer.Dial(network, addr)
			},
			DisableKeepAlives: true,
		},
		Timeout: 30 * time.Second,
	}, torDialer, nil
}

// func newEthereumWithTor(name string, rpcUrl string, wsRpcUrl string) (*ethclient.Client, error) {
// 	// Create Tor HTTP client
// 	torClient, err := NewTorHTTPClient()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create custom RPC client with Tor
// 	rpcClient, err := rpc.DialHTTPWithClient(rpcUrl, torClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ethClient := ethclient.NewClient(rpcClient)

// 	wsRPCClient, err := rpc.DialWebsocket(context.Background(), wsRpcUrl, "", torClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	wsEthClient := ethclient.NewClient(wsRPCClient)

// 	return ethClient, nil
// }
