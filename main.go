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

	blockNumber, err := eth.GetBlockNumber()

	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println(blockNumber)
	fmt.Println(client.BlockByNumber(context.TODO(), blockNumber))
	fmt.Println(client.HeaderByNumber(context.TODO(), big.NewInt(0)))
}
