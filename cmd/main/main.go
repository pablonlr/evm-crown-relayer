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

	err = relayer.Run(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	block := make(chan struct{})
	<-block

}
