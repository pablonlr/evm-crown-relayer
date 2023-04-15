package jobs

import (
	crownd "github.com/pablonlr/go-rpc-crownd"
	"github.com/pablonlr/poly-crown-relayer/crown"
	rtypes "github.com/pablonlr/poly-crown-relayer/types"
)

const Protocol_Confirmations_Required = 6

type CrownProtocolJob struct {
	protocolID                 string
	protocolOwner              string
	protocolName               string
	protocolRegTx              string
	description                string
	protocolRegTxConfirmations int
	crwResolver                *crown.CrownResolver
}

func NewCrownProtocolJob(protocolID, protocolName, protocolOwner, protocolDescription string, crwResolver *crown.CrownResolver) *CrownProtocolJob {
	return &CrownProtocolJob{
		protocolID:                 protocolID,
		protocolName:               protocolName,
		protocolOwner:              protocolOwner,
		description:                protocolDescription,
		protocolRegTx:              "",
		protocolRegTxConfirmations: 0,
		crwResolver:                crwResolver,
	}
}

func (c *CrownProtocolJob) GetNextTask(previousTask *rtypes.Task) *rtypes.Task {
	if previousTask == nil {
		return &rtypes.Task{
			ID:         GetNftProtoTask,
			Exec:       c.crwResolver.GetNFTProtocol,
			ExecParams: []string{c.protocolID},
		}
	}

	switch previousTask.ID {
	case GetNftProtoTask:
		if previousTask.TResult.ResultValue == nil && previousTask.TResult.Err.Err == nil {
			return &rtypes.Task{
				ID:         RegisterNftProtoTask,
				Exec:       c.crwResolver.RegisterNFTProtocol,
				ExecParams: []string{c.protocolID, c.protocolName, c.protocolOwner, c.description},
			}
		}
		proto, ok := previousTask.TResult.ResultValue.(*crownd.NFTProtocol)
		if !ok {
			return nil
		}
		c.protocolRegTx = proto.RegistrationTxHash
		return &rtypes.Task{
			ID:         NftProtocolConfirmedTask,
			Exec:       c.crwResolver.NFTProtocolConfirmed,
			ExecParams: []string{c.protocolID},
		}
	case RegisterNftProtoTask:
		if previousTask.TResult.Err.Err != nil {
			return nil
		}
		txId, ok := previousTask.TResult.ResultValue.(string)
		if !ok {
			return nil
		}
		c.protocolRegTx = txId
		return &rtypes.Task{
			ID:         WaitConfirmationsTask,
			Exec:       WaitTime,
			ExecParams: []string{"360"},
		}
	case WaitConfirmationsTask:
		return &rtypes.Task{
			ID:         NftProtocolConfirmedTask,
			Exec:       c.crwResolver.NFTProtocolConfirmed,
			ExecParams: []string{c.protocolID},
		}
	case NftProtocolConfirmedTask:
		if previousTask.TResult.Err.Err != nil {
			return nil
		}
		confirmations, ok := previousTask.TResult.ResultValue.(int)
		if !ok {
			return nil
		}

		if confirmations < Protocol_Confirmations_Required {
			return &rtypes.Task{
				ID:         WaitConfirmationsTask,
				Exec:       WaitTime,
				ExecParams: []string{"360"},
			}
		}
		c.protocolRegTxConfirmations = confirmations
		return nil
	default:
		return nil
	}

}

func (c *CrownProtocolJob) TaskCount() int {
	return 3
}

func (c *CrownProtocolJob) Result() *rtypes.JobResult {
	if c.protocolRegTx != "" && c.protocolRegTxConfirmations > Protocol_Confirmations_Required {
		return &rtypes.JobResult{
			Finished: true,
			Value:    c.protocolRegTx,
		}
	}
	return &rtypes.JobResult{
		Finished: false,
	}
}
