package jobs

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	crownd "github.com/pablonlr/go-rpc-crownd"
	"github.com/pablonlr/poly-crown-relayer/crown"
	rtypes "github.com/pablonlr/poly-crown-relayer/types"
)

type CrownRegEvent struct {
	TokenID    *big.Int
	Uri        string
	CrwAddress string
}

type CrownRegistrationJob struct {
	event           *CrownRegEvent
	crownTokenID    string
	crwResolver     *crown.CrownResolver
	crownProtocolId string
	txId            string
	txConfirmations int
	invalidAddress  bool
}

func NewCrownRegistrationJob(eventName string, cabi *abi.ABI, vLog *types.Log, crwResolver *crown.CrownResolver, crownProtocol string) (*CrownRegistrationJob, error) {
	evnt := &CrownRegEvent{}
	err := cabi.UnpackIntoInterface(evnt, "CrownRegistration", vLog.Data)
	if err != nil {
		return nil, err
	}
	return &CrownRegistrationJob{
		event:           evnt,
		crwResolver:     crwResolver,
		crownTokenID:    getCrownID(evnt.TokenID, evnt.CrwAddress, vLog.TxHash.String()),
		crownProtocolId: crownProtocol,
	}, nil
}

func getCrownID(tokendId *big.Int, crwAddress, evmTxId string) string {
	st := tokendId.String()
	st += crwAddress
	st += evmTxId
	return hash(st)
}

func (c *CrownRegistrationJob) GetNextTask(previousTask *rtypes.Task) *rtypes.Task {

	if previousTask == nil {
		return &rtypes.Task{
			ID:         GetNFTokenTask,
			Exec:       c.crwResolver.GetNFToken,
			ExecParams: []string{c.crownProtocolId, c.crownTokenID},
		}
	}

	switch previousTask.ID {
	case GetNFTokenTask:
		if previousTask.TResult.ResultValue == nil && previousTask.TResult.Err.Err == nil {
			return &rtypes.Task{
				ID:         RegisterNFToken,
				Exec:       c.crwResolver.RegisterNFToken,
				ExecParams: []string{c.crownProtocolId, c.crownTokenID, c.event.CrwAddress, c.event.Uri},
			}
		}
		nfToken, ok := previousTask.TResult.ResultValue.(*crownd.NFToken)
		if !ok {
			return nil
		}
		if len(nfToken.BlockHash) == 64 {
			c.txId = nfToken.RegistrationTxHash
			c.txConfirmations = 1
			return nil
		}
		return &rtypes.Task{
			ID:         WaitConfirmationsTask,
			Exec:       WaitTime,
			ExecParams: []string{"60"},
		}
	case RegisterNFToken:
		if previousTask.TResult.Err.Err != nil {
			if previousTask.TResult.Err.Code == rtypes.InvalidCrownAddress {
				log.Println(previousTask.TResult.Err.Err.Error(), c.event.CrwAddress)
				c.invalidAddress = true
			}
			return nil
		}

		txId, ok := previousTask.TResult.ResultValue.(string)
		if !ok {
			return nil
		}
		c.txId = txId
		return &rtypes.Task{
			ID:         WaitConfirmationsTask,
			Exec:       WaitTime,
			ExecParams: []string{"60"},
		}
	case WaitConfirmationsTask:
		txId := strings.Replace(c.txId, "\"", "", -1)
		return &rtypes.Task{
			ID:         IsConfirmedNftTx,
			Exec:       c.crwResolver.NFTokenConfirmed,
			ExecParams: []string{txId},
		}
	case IsConfirmedNftTx:
		if previousTask.TResult.Err.Err != nil && strings.HasPrefix("Error -1 : Can't find an NFT record by tx id", previousTask.TResult.Err.Err.Error()) {
			return &rtypes.Task{
				ID:         WaitConfirmationsTask,
				Exec:       WaitTime,
				ExecParams: []string{"60"},
			}
		}
		if previousTask.TResult.Err.Err != nil {
			log.Println(previousTask.TResult.Err.Err.Error())
			return nil
		}
		confs, ok := previousTask.TResult.ResultValue.(int)
		if !ok {
			return nil
		}
		if confs > 0 {
			c.txConfirmations = confs
			return nil
		}
		return &rtypes.Task{
			ID:         WaitConfirmationsTask,
			Exec:       WaitTime,
			ExecParams: []string{"60"},
		}
	}
	return nil
}

func (c *CrownRegistrationJob) TaskCount() int {
	return 3
}

func (c *CrownRegistrationJob) Result() *rtypes.JobResult {
	if c.invalidAddress {
		return &rtypes.JobResult{
			Finished: true,
			Value:    "invalid address",
		}
	}
	if c.txId != "" && c.txConfirmations > 0 {
		return &rtypes.JobResult{
			Finished: true,
			Value:    c.txId,
		}
	}
	return &rtypes.JobResult{
		Finished: false,
	}

}

func hash(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
