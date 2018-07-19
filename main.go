package main

import (
	"context"
	"fmt"
	"go-eth/eth"
	"math/big"
)

func main() {
	client, err := eth.Connect("http://localhost:8545")

	if err != nil {
		fmt.Errorf(err.Error())
	}

	blockNumber, err := client.GetBlockNumber()

	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println(blockNumber)
	fmt.Println(client.EthClient.BlockByNumber(context.TODO(), blockNumber))
	fmt.Println(client.EthClient.HeaderByNumber(context.TODO(), big.NewInt(0)))
}
