package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"go-eth/contract"
	"go-eth/eth"
	"os"
)

var (
	errEmptyEtherscanToken     = fmt.Errorf("should set ETHERSCAN_TOKEN to environment")
	errEmptyUniswapPairAddress = fmt.Errorf("should set UNISWAP_ADDRESS to environment")
)

func printInfo(msg string) {
	str := fmt.Sprintf("[INFO] %s", msg)
	fmt.Println(str)
}

func printError(err error) {
	str := fmt.Sprintf("[ERROR] %s", err.Error())
	fmt.Println(str)
}

func main() {
	rpc := os.Getenv("rpc")
	if len(rpc) == 0 {
		rpc = "http://localhost:8545"
	}
	escToken := os.Getenv("ETHERSCAN_TOKEN")
	if len(escToken) == 0 {
		printError(errEmptyEtherscanToken)
		return
	}
	pairAddr := os.Getenv("UNISWAP_ADDRESS")
	if len(pairAddr) == 0 {
		printError(errEmptyUniswapPairAddress)
		return
	}
	dPairAddr := common.HexToAddress(pairAddr)
	escClient := contract.NewEtherscan(escToken)
	strABI, err := escClient.GetVerifiedABI(context.TODO(), dPairAddr)
	if err != nil {
		printError(err)
		return
	}
	if client, err := eth.Connect(rpc); err != nil {
		printError(err)
		return
	} else {
		chainID, err := client.ChainID(context.TODO())
		if err != nil {
			printError(err)
			return
		}
		printInfo(fmt.Sprintf("Connect to chain %d", chainID.Int64()))
		contract, err := contract.NewContractWithABI(strABI, dPairAddr, client.EthClient)
		if err != nil {
			printError(err)
			return
		}
		fmt.Println("Got uniswap contract instance: ", contract)
		return
	}
}
