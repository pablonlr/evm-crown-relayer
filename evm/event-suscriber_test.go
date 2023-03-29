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

const testABI = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"tokenID","type":"uint256"},{"indexed":false,"internalType":"string","name":"uri","type":"string"},{"indexed":false,"internalType":"string","name":"crwAddress","type":"string"}],"name":"CrownRegistration","type":"event"}]`

func getRPC() (rpc string) {
	rpc = os.Getenv("TEST_RPC")
	return
}

func TestGetLogsToBlockN(t *testing.T) {
	contractAddress := "0xCB7d76b8C525C4EefC95db4c9D30BeE1A401C1DA"
	polygonRPC := getRPC()
	if len(polygonRPC) < 1 {
		log.Fatalf("Please set up enviroment variable TEST_RPC")
	}
	log.Println(polygonRPC)
	eventHex := "0xdaa2dbe0deba8671ce45936bb34dae41efc232dcccb282c8b02f42500ecbf432"
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
	logs, err := suscrib.GetLogsToBlockN(currentH)
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
	logs, err := suscrib.GetLogsToBlockN(currentH)
	if err != nil {
		panic(err)
	}
	if len(logs) < 1 {
		t.Errorf("No logs fetched")
	}

}
