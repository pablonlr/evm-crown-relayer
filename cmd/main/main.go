package main

import (
	"context"
	"log"

	"github.com/pablonlr/poly-crown-relayer/config"
	"github.com/pablonlr/poly-crown-relayer/relayer"
)

func main() {
	conf, err := config.LoadConfig("./config.json")
	if err != nil {
		panic(err)
	}

	relayer, err := relayer.NewRelayerFromConf(*conf)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = relayer.Run(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for {
	}

}
