package relayer

import (
	"context"
	"fmt"
	"log"

	"github.com/pablonlr/poly-crown-relayer/config"
	"github.com/pablonlr/poly-crown-relayer/crown"
	"github.com/pablonlr/poly-crown-relayer/evm"
	"github.com/pablonlr/poly-crown-relayer/relayer/jobs"
	rtypes "github.com/pablonlr/poly-crown-relayer/types"
)

type Instance struct {
	Name         string
	evmSuscriber *evm.Suscriber
	contract     evm.Contract
	crwResolver  *crown.CrownResolver
	worker       *Worker
	protoDetails *config.CrownRegConfig
	protocolJob  *rtypes.Job
}

func NewInstance(name string, suscriber *evm.Suscriber, crwResolver *crown.CrownResolver, regConf *config.CrownRegConfig) *Instance {

	return &Instance{
		Name:         name,
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
	logsChan, errChan, err := i.evmSuscriber.GetPastLogsAndSuscribeToFutureLogsRest(ctx)
	if err != nil {
		return err
	}
	go i.worker.Start(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping instance %s...", i.Name)
			return nil
		case err := <-errChan:
			return err
		case lg := <-logsChan:
			taskBuilder, err := jobs.NewCrownRegistrationJob(i.evmSuscriber.GetEventName(), i.contract.GetAbi(), &lg, i.crwResolver, i.protoDetails.ProtocolID)
			if err != nil {
				return err
			}
			job := rtypes.NewJob(jobs.RegisterNFTokenJob, taskBuilder)

			select {
			case i.worker.inputChan <- job:
			default:
				log.Printf("Failed to send job to worker in instance %s: worker input channel closed", i.Name)
				return fmt.Errorf("worker input channel closed")
			}
		}

	}

}
