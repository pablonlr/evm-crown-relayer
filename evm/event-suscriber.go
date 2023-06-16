package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pablonlr/poly-crown-relayer/config"
)

type Suscriber struct {
	eventHash        common.Hash
	eventName        string
	indexedFromBlock *big.Int
	indexedToBlock   *big.Int
	contract         *Contract
	resolver         *EVMResolver
}

func NewSuscriber(contract *Contract, resolver *EVMResolver, eventHashHex, eventName string, contractStart *big.Int) *Suscriber {

	eventHash := common.HexToHash(eventHashHex)
	return &Suscriber{
		eventHash:        eventHash,
		indexedFromBlock: contractStart,
		contract:         contract,
		resolver:         resolver,
		eventName:        eventName,
	}
}

func NewSuscriberFromConf(solDef config.SolDefinitions, evmConfig config.EVMconfig) (*Suscriber, error) {

	contract, err := NewContract(evmConfig.ChainName, evmConfig.ContractAddress, solDef.ContractABI)
	if err != nil {
		return nil, err
	}
	indexedFrom := big.NewInt(evmConfig.SuscribeFromBlock)
	resolver, err := NewEVMResolver(evmConfig.RPC)
	if err != nil {
		return nil, err
	}

	eventHash := common.HexToHash(solDef.EventHex)
	return &Suscriber{
		eventHash:        eventHash,
		indexedFromBlock: indexedFrom,
		contract:         contract,
		resolver:         resolver,
		eventName:        solDef.EventName,
	}, nil
}

func (s *Suscriber) GetLogsToBlockN(toBlock *big.Int) ([]types.Log, error) {
	query := s.getNewQuery(toBlock)
	return s.resolver.client.FilterLogs(context.Background(), query)
}

func (s *Suscriber) GetPastLogsAndSuscribeToFutureLogs(ctx context.Context) (chan types.Log, <-chan error, error) {
	query := s.getNewQuery(nil)
	out := make(chan types.Log)
	currentH, err := s.resolver.CurrentBlockHeight()
	if err != nil {
		return nil, nil, err
	}
	logs, err := s.GetLogsToBlockN(currentH)
	if err != nil {
		return nil, nil, err
	}

	//sync channel to send the signal when the past logs are sent
	syncCh := make(chan struct{})

	// Send past logs to the channel
	go func() {
		for _, log := range logs {
			out <- log
		}

		close(syncCh)
	}()

	// Subscribe to future logs
	futureLogs := make(chan types.Log)
	sub, err := s.resolver.client.SubscribeFilterLogs(ctx, query, futureLogs)
	if err != nil {
		return nil, nil, err
	}

	errChan := make(chan error)
	go func() {

		defer close(errChan)
		defer close(out)

		// wait for the past logs to be sent
		<-syncCh

		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case log, ok := <-futureLogs:
				if ok {
					out <- log
				} else {
					return
				}
			case err := <-sub.Err():
				errChan <- fmt.Errorf("error while subscribing to future logs: %v", err)
				return
			}
		}
	}()

	return out, errChan, nil
}

func (s *Suscriber) getNewQuery(to *big.Int) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: s.indexedFromBlock,
		ToBlock:   to,
		Addresses: []common.Address{
			s.contract.contractAddress,
		},
		Topics: [][]common.Hash{{
			s.eventHash,
		}},
	}
}

func (s *Suscriber) GetContract() Contract {
	return *s.contract
}

func (s *Suscriber) GetEventName() string {
	return s.eventName
}
