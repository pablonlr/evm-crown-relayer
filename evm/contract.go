package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	ChainID         string
	contractAddress common.Address
	tokenName       string
	tokenSymbol     string
	abi             *abi.ABI
}

func NewContract(chainID, address string, contractAbi *abi.ABI) (*Contract, error) {
	contractAddress := common.HexToAddress(address)
	/*
		contractAbi, err := abi.JSON(strings.NewReader(string(cABI)))
		if err != nil {
			log.Fatal(err)
		}
	*/

	return &Contract{
		ChainID:         chainID,
		contractAddress: contractAddress,
		abi:             contractAbi,
	}, nil
}

func (c *Contract) GetAbi() *abi.ABI {
	return c.abi
}
