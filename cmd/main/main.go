package main

import (
	"context"
	"log"

	"github.com/pablonlr/poly-crown-relayer/config"
	"github.com/pablonlr/poly-crown-relayer/crown"
	"github.com/pablonlr/poly-crown-relayer/evm"
	"github.com/pablonlr/poly-crown-relayer/relayer"
)

func main() {
	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		panic(err)
	}
	crwResolver, err := crown.NewCrownResolver(*conf.CrownClientConf)
	if err != nil {
		panic(err)
	}

	suscrib, err := evm.NewSuscriberFromConf(*conf.Definitions, *conf.Instances[0].EVM)
	if err != nil {
		panic(err)
	}

	instanc := relayer.NewInstance(suscrib, crwResolver, conf.Instances[0].Crown)
	err = instanc.ConfigureProtocol()
	if err != nil {
		panic(err)
	}
	log.Println("Protocol configured")
	err = instanc.StartRegistrations(context.Background())
	if err != nil {
		panic(err)
	}

}
