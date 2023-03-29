package config

import (
	"fmt"
	"strings"
)

func onlyContains(s string, setOfChars string) bool {
	for _, c := range s {
		if !strings.ContainsRune(setOfChars, c) {
			return false
		}
	}
	return true
}

func formatConfig(config *Config) error {
	for i, inst := range config.Instances {
		if len(inst.Name) == 0 {
			return fmt.Errorf("Instance %d must have a name", i)
		}
		if len(inst.EVM.ChainName) == 0 {
			return fmt.Errorf("Instance %d must have a chain name", i)
		}
		if len(inst.EVM.ContractAddress) == 0 {
			return fmt.Errorf("Instance %d must have a contract address", i)
		}
		if len(inst.Crown.ProtocolID) > 12 {
			return fmt.Errorf("Instance %d must have a valid Crown ProtocolID", i)
		}
		if len(inst.Crown.ProtocolName) > 24 {
			return fmt.Errorf("Instance %d must have a valid Crown Protocol Name", i)
		}
		if len(inst.EVM.RPC) == 0 {
			return fmt.Errorf("Instance %d must have a valid RPC", i)
		}
		if inst.EVM.SuscribeFromBlock < 0 {
			return fmt.Errorf("Instance %d must have a suscribe from block greater or equal than 0", i)
		}
		if len(inst.EVM.ContractAddress) != 42 {
			return fmt.Errorf("Instance %d contract address must be 42 characters long", i)
		}
		if len(inst.Crown.ProtocolOwnerAddress) != 36 {
			return fmt.Errorf("Instance %d protocol owner address must be 36 characters long", i)
		}
		if !onlyContains(inst.Crown.ProtocolID, ".abcdefghijklmnopqrstuvwxyz12345") {
			return fmt.Errorf("Instance %d protocol ID must only contain lowercase letters and numbers", i)
		}
		if len(inst.Crown.ProtocolDescription) == 0 {
			inst.Crown.ProtocolDescription = fmt.Sprintf("%s contract from %s to Crown", inst.EVM.ContractAddress, inst.EVM.ChainName)
		}
	}
	return nil
}
