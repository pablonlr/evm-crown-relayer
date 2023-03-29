package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	RPC_USER   = "RPC_USER"
	RPC_PASS   = "RPC_PASS"
	UnlockPass = "UNLOCK_PASS"
)

type SolDefinitions struct {
	ContractABI *abi.ABI `json:"contract_abi"`
	EventName   string   `json:"event_name"`
	EventHex    string   `json:"event_hex"`
}

type EVMconfig struct {
	ChainName         string `json:"chain_name"`
	RPC               string `json:"rpc"`
	ContractAddress   string `json:"contract_address"`
	SuscribeFromBlock int64  `json:"suscribe_from_block"`
}

type CrownClientConfig struct {
	ClientAdress         string `json:"client_address"`
	ClientPort           int    `json:"client_port"`
	ClientRequestTimeout int    `json:"client_request_timeout"`
	UnlockNSeconds       int    `json:"unlock_n_seconds"`
	Secrets              *CrownSecrets
}

type CrownRegConfig struct {
	ProtocolID           string `json:"protocol_id"`
	ProtocolName         string `json:"protocol_name"`
	ProtocolOwnerAddress string `json:"protocol_owner_address"`
	ProtocolDescription  string `json:"protocol_description"`
}

type RelayerInstance struct {
	Name     string          `json:"name"`
	EVM      *EVMconfig      `json:"evm_config"`
	Crown    *CrownRegConfig `json:"crown_protocol_config"`
	RetryJob int             `json:"retry_job"`
}

type Config struct {
	Instances       []RelayerInstance  `json:"instances"`
	CrownClientConf *CrownClientConfig `json:"crown_client_config"`
	Definitions     *SolDefinitions    `json:"sol_definitions"`
}

type CrownSecrets struct {
	RPC_USER   string
	RPC_PASS   string
	UnlockPass string
}

func LoadConfig(path string) (*Config, error) {
	config, err := LoadPublicConfig(path)
	if err != nil {
		return nil, err
	}
	err = formatConfig(config)
	if err != nil {
		return nil, err
	}
	if config.CrownClientConf.Secrets != nil {
		return config, nil // secrets already loaded from JSON
	}

	ruser := os.Getenv(RPC_USER)
	if ruser == "" {
		return nil, errors.New("error reading enviroment var RPC_USER")
	}

	rpass := os.Getenv(RPC_PASS)
	if rpass == "" {
		return nil, errors.New("error reading enviroment var RPC_PASS")
	}
	pass := os.Getenv(UnlockPass)
	if pass == "" {
		return nil, errors.New("error reading enviroment var UnlockPass")
	}
	sec := CrownSecrets{
		RPC_USER:   ruser,
		RPC_PASS:   rpass,
		UnlockPass: pass,
	}
	config.CrownClientConf.Secrets = &sec
	return config, nil
}

func LoadPublicConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err

	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	if len(config.Instances) == 0 {
		return nil, errors.New("no instances found in config file")
	}

	return config, err
}
