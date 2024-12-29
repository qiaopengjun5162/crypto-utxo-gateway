package base

import (
	"github.com/ethereum/go-ethereum/log"

	"github.com/btcsuite/btcd/rpcclient"
)

type Client struct {
	*rpcclient.Client
	compressed bool
}

func NewBaseClient(RpcUrl, RpcUser, RpcPass string) (*Client, error) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         RpcUrl,
		User:         RpcUser,
		Pass:         RpcPass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}, nil)
	if err != nil {
		log.Error("new bitcoin rpc client fail", "err", err)
		return nil, err
	}
	return &Client{
		Client:     client,
		compressed: true,
	}, nil
}
