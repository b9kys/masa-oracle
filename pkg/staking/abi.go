// pkg/staking/abi.go
package staking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	MasaTokenABIPath         = "contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/MasaToken.sol/MasaToken.json"
	OracleNodeStakingABIPath = "contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/OracleNodeStaking.sol/OracleNodeStaking.json"
)

type ContractAddresses struct {
	Sepolia struct {
		MasaToken         string `json:"MasaToken"`
		OracleNodeStaking string `json:"OracleNodeStaking"`
		StakedMasaToken   string `json:"StakedMasaToken"`
	} `json:"sepolia"`
}

func GetABI(jsonPath string) (abi.ABI, error) {
	jsonFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI: %v", err)
	}

	var contract struct {
		ABI json.RawMessage `json:"abi"`
	}
	err = json.Unmarshal(jsonFile, &contract)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to unmarshal ABI JSON: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contract.ABI)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return parsedABI, nil
}
