package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EVMResolver struct {
	client *ethclient.Client
}

func NewEVMResolver(rpc string) (*EVMResolver, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	return &EVMResolver{
		client: client,
	}, nil

}

func (r *EVMResolver) CurrentBlockHeight() (*big.Int, error) {
	header, err := r.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}
