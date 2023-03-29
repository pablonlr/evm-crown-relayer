package relayer
/*
import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

func (i *CrownRegistrationInterpeter) logToCrown(vLog *types.Log) (*CrownRegEvent, error) {
	evnt := &CrownRegEvent{}
	err := i.contractAbi.UnpackIntoInterface(evnt, "CrownRegistration", vLog.Data)
	if err != nil {
		return nil, err
	}
	return evnt, nil
}

type CrownRegistrationInterpeter struct {
	eventName   string
	contractAbi *abi.ABI
}

func NewCrownRegistrationInterpeter(eventName string, cabi *abi.ABI) CrownRegistrationInterpeter {
	return CrownRegistrationInterpeter{
		eventName:   eventName,
		contractAbi: cabi,
	}
}
*/
