package evm

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pablonlr/poly-crown-relayer/config"
)

const RestFetchWait = 120

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

func (s *Suscriber) GetLogsFromBlockMToBlockN(fromBlock, toBlock *big.Int) ([]types.Log, error) {
	query := s.getNewQuery(fromBlock, toBlock)
	return s.resolver.client.FilterLogs(context.Background(), query)
}

// Rest implementation of suscribe filter logs query
func (s *Suscriber) GetPastLogsAndSuscribeToFutureLogsRest(ctx context.Context) (chan types.Log, <-chan error, error) {
	currentH, err := s.resolver.CurrentBlockHeight()
	if err != nil {
		return nil, nil, err
	}
	out := make(chan types.Log)
	logs, err := s.GetLogsFromBlockMToBlockN(s.indexedFromBlock, currentH)
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

	errChan := make(chan error)
	go func() {
		defer close(errChan)
		defer close(out)

		<-syncCh

		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case <-time.After(RestFetchWait * time.Second):
				// Fetch latest block height
				newH, err := s.resolver.CurrentBlockHeight()
				if err != nil {
					errChan <- err
					return
				}
				if newH.Cmp(currentH) > 0 {
					nextBlock := new(big.Int).Add(currentH, big.NewInt(1))
					newLogs, err := s.GetLogsFromBlockMToBlockN(nextBlock, newH)
					if err != nil {
						errChan <- err
						return
					}
					for _, log := range newLogs {
						out <- log
					}
					currentH = newH
				}
			}
		}
	}()

	return out, errChan, nil

}

func (s *Suscriber) GetPastLogsAndSuscribeToFutureLogs(ctx context.Context) (chan types.Log, <-chan error, error) {
	query := s.getNewQuery(s.indexedFromBlock, nil)
	out := make(chan types.Log)
	currentH, err := s.resolver.CurrentBlockHeight()
	if err != nil {
		return nil, nil, err
	}
	logs, err := s.GetLogsFromBlockMToBlockN(s.indexedFromBlock, currentH)
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

	futureLogs := make(chan types.Log)
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		defer close(out)

		<-syncCh

		for {
			sub, err := s.resolver.client.SubscribeFilterLogs(ctx, query, futureLogs)
			if err != nil {
				log.Printf("error al suscribirse a los logs futuros: %v", err)
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}
			defer sub.Unsubscribe()

			for {
				select {
				case err := <-sub.Err():
					log.Printf("error en la suscripción: %v", err)
					select {
					case <-ctx.Done():
						errChan <- ctx.Err()
						return
					case <-time.After(5 * time.Second):
						break
					}
				case log := <-futureLogs:
					out <- log
				}
			}
		}
	}()

	return out, errChan, nil
}

func (s *Suscriber) getNewQuery(from *big.Int, to *big.Int) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: from,
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
