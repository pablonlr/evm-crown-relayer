package relayer

import (
	"context"
	"fmt"

	"github.com/pablonlr/poly-crown-relayer/config"
	"github.com/pablonlr/poly-crown-relayer/crown"
	"github.com/pablonlr/poly-crown-relayer/evm"
	"github.com/pablonlr/poly-crown-relayer/relayer/jobs"
	rtypes "github.com/pablonlr/poly-crown-relayer/types"
)

type Instance struct {
	evmSuscriber *evm.Suscriber
	contract     evm.Contract
	crwResolver  *crown.CrownResolver
	worker       *Worker
	protoDetails *config.CrownRegConfig
	protocolJob  *rtypes.Job
}

func NewInstance(suscriber *evm.Suscriber, crwResolver *crown.CrownResolver, regConf *config.CrownRegConfig) *Instance {

	return &Instance{
		evmSuscriber: suscriber,
		contract:     suscriber.GetContract(),
		crwResolver:  crwResolver,
		worker:       NewWorker(),
		protoDetails: regConf,
	}
}

func (i *Instance) ConfigureProtocol() error {
	taskBuilder := jobs.NewCrownProtocolJob(i.protoDetails.ProtocolID, i.protoDetails.ProtocolName, i.protoDetails.ProtocolOwnerAddress, i.protoDetails.ProtocolDescription, i.crwResolver)
	job := rtypes.NewJob(jobs.NftProtoJob, taskBuilder)
	i.protocolJob = job
	return i.worker.ProcessJob(job)
}

func (i *Instance) StartRegistrations(ctx context.Context) error {
	if i.protocolJob == nil {
		return fmt.Errorf("protocol job not configured")
	}
	if !i.protocolJob.Tasks.Result().Finished {
		return fmt.Errorf("protocol job not completed")
	}
	logsChan, errChan, err := i.evmSuscriber.GetPastLogsAndSuscribeToFutureLogs()
	if err != nil {
		return err
	}
	go i.worker.Start(ctx)

	for {
		select {
		case err := <-errChan:
			return err
		case log := <-logsChan:
			taskBuilder, err := jobs.NewCrownRegistrationJob(i.evmSuscriber.GetEventName(), i.contract.GetAbi(), &log, i.crwResolver, i.protoDetails.ProtocolID)
			if err != nil {
				return err
			}
			job := rtypes.NewJob(jobs.RegisterNFTokenJob, taskBuilder)
			i.worker.inputChan <- job
		}
	}

}
