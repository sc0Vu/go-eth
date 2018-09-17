package main

import (
	"context"
	"fmt"
	"go-eth/eth"
	"math/big"
	// "time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	client, err := eth.Connect("http://localhost:8545")

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	blockNumber, err := client.GetBlockNumber(context.TODO())

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	fmt.Printf("Start to call eth_sendTransaction\nLatest block number: %s\n", blockNumber.String())
	
	from := common.HexToAddress("30b82c8694b59695d78f33a7ba1c2a55dfa618d5")
	to := common.HexToAddress("5e0f92917d632f7cdb7564a67644ca45430b524c")
	amount := big.NewInt(1)
	gasLimit := big.NewInt(90000)
	gasPrice := big.NewInt(0)
	data := []byte{}
	message := eth.NewMessage(&from, &to, amount, gasLimit, gasPrice, data)

	if txHash, err := client.SendTransaction(context.TODO(), &message); err != nil {
		fmt.Errorf(err.Error())
		return
	} else {
		fmt.Printf("Message: %s\nTransaction has been sent, transaction hash: %s\n", message.String(), txHash.String())
	    tx, isPending, _ := client.EthClient.TransactionByHash(context.TODO(), txHash)
		fmt.Printf("Transaction nonce: %d\nTransaction pending: %v\n", tx.Nonce(), isPending)
		// check transaction receipt
		receiptChan := make(chan *types.Receipt)
		client.CheckTransaction(context.TODO(), receiptChan, txHash, 1)
		receipt := <-receiptChan
		fmt.Printf("Transaction status: %v\n", receipt.Status)
	}
}
