package evm

import (
	"log"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pablonlr/poly-crown-relayer/config"
)

const testABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"emisor","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenID","type":"uint256"},{"indexed":false,"internalType":"string","name":"uri","type":"string"},{"indexed":false,"internalType":"string","name":"crwAddress","type":"string"}],"name":"CrownRegistration","type":"event"}]`

func getRPC() (rpc string) {
	rpc = os.Getenv("TEST_RPC")
	return
}

func TestGetLogsToBlockN(t *testing.T) {
	contractAddress := "0xE856B7Bdc52F86e97f6197A825f505f17046589d"
	polygonRPC := getRPC()
	if len(polygonRPC) < 1 {
		log.Fatalf("Please set up enviroment variable TEST_RPC")
	}
	log.Println(polygonRPC)
	eventHex := "0x32abe70a521c5c3431eb1ac5dfd23738bae8e0a25afde8918617300bb9c379ca"
	contractAbi, err := abi.JSON(strings.NewReader(string(testABI)))
	if err != nil {
		log.Fatal(err)
	}
	contract, err := NewContract("PolygonMumbai", contractAddress, &contractAbi)
	if err != nil {
		panic(err)
	}
	resolver, err := NewEVMResolver(polygonRPC)
	if err != nil {
		panic(err)
	}
	suscrib := NewSuscriber(contract, resolver, eventHex, "CrownRegistration", big.NewInt(0))
	currentH, err := suscrib.resolver.CurrentBlockHeight()
	if err != nil {
		panic(err)
	}
	logs, err := suscrib.GetLogsFromBlockMToBlockN(suscrib.indexedFromBlock, currentH)
	if err != nil {
		panic(err)
	}
	if len(logs) < 1 {
		t.Errorf("No logs fetched")
	}
}

func TestGetLogsToBlockNFromConfFile(t *testing.T) {
	conf, err := config.LoadConfig("config_test.json")
	if err != nil {
		panic(err)
	}
	suscrib, err := NewSuscriberFromConf(*conf.Definitions, *conf.Instances[0].EVM)
	if err != nil {
		panic(err)
	}
	currentH, err := suscrib.resolver.CurrentBlockHeight()
	if err != nil {
		panic(err)
	}
	logs, err := suscrib.GetLogsFromBlockMToBlockN(suscrib.indexedFromBlock, currentH)
	if err != nil {
		panic(err)
	}
	if len(logs) < 1 {
		t.Errorf("No logs fetched")
	}

}
func TestGetLogsPastAndFuturesFromConfFile(t *testing.T) {
	conf, err := config.LoadConfig("config_test.json")
	if err != nil {
		panic(err)
	}
	suscrib, err := NewSuscriberFromConf(*conf.Definitions, *conf.Instances[0].EVM)
	if err != nil {
		panic(err)
	}
	currentH, err := suscrib.resolver.CurrentBlockHeight()
	if err != nil {
		panic(err)
	}
	logs, err := suscrib.GetLogsFromBlockMToBlockN(suscrib.indexedFromBlock, currentH)
	if err != nil {
		panic(err)
	}
	if len(logs) < 1 {
		t.Errorf("No logs fetched")
	}

}
