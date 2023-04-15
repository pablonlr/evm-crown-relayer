package relayer

import (
	"context"
	"log"

	"github.com/pablonlr/poly-crown-relayer/config"
	"github.com/pablonlr/poly-crown-relayer/crown"
	"github.com/pablonlr/poly-crown-relayer/evm"
)

type Relayer struct {
	Instances []Instance
}

func NewRelayer() *Relayer {
	return &Relayer{}
}

func (r *Relayer) AddInstance(instance Instance) {
	r.Instances = append(r.Instances, instance)
}

func startInstance(i *Instance, ctx context.Context) error {
	err := i.ConfigureProtocol()
	if err != nil {
		log.Printf("Error configuring protocol for instance %s, stopping this instance...", i.Name)
		return err
	}
	log.Println("Protocol configured for instance: ", i.Name)
	err = i.StartRegistrations(ctx)
	if err != nil {
		log.Printf("Error in execution for instance %s, stopping this instance...", i.Name)
		return err
	}
	return nil

}

func (r *Relayer) Run(ctx context.Context) error {
	for _, instance := range r.Instances {
		go startInstance(&instance, ctx)
	}
	return nil
}

func NewRelayerFromConf(conf config.Config) (*Relayer, error) {
	relayer := NewRelayer()
	for _, instanceConf := range conf.Instances {
		suscriber, err := evm.NewSuscriberFromConf(*conf.Definitions, *instanceConf.EVM)
		if err != nil {
			return nil, err
		}
		crwResolver, err := crown.NewCrownResolver(*conf.CrownClientConf)
		if err != nil {
			return nil, err
		}

		instance := NewInstance(instanceConf.Name, suscriber, crwResolver, instanceConf.Crown)
		relayer.AddInstance(*instance)
	}
	return relayer, nil
}
