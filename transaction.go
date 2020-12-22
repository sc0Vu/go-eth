package main

import (
	"context"
	"fmt"
	"go-eth/eth"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	client, err := eth.Connect(rpc)
	if err != nil {
		printError(err)
		return
	}

	blockNumber, err := client.GetBlockNumber(context.TODO())
	if err != nil {
		printError(err)
		return
	}
	printInfo("Start to call eth_sendTransaction")
	printInfo(fmt.Sprintf("Latest block number: %d", blockNumber))

	from := common.HexToAddress("30b82c8694b59695d78f33a7ba1c2a55dfa618d5")
	to := common.HexToAddress("5e0f92917d632f7cdb7564a67644ca45430b524c")
	amount := big.NewInt(1)
	gasLimit := big.NewInt(90000)
	gasPrice := big.NewInt(0)
	data := []byte{}
	message := eth.NewMessage(&from, &to, amount, gasLimit, gasPrice, data)

	if txHash, err := client.SendTransaction(context.TODO(), &message); err != nil {
		printError(err)
		return
	} else {
		printInfo(fmt.Sprintf("Message: %s Transaction has been sent, transaction hash: %s", message.String(), txHash.String()))
		tx, isPending, _ := client.EthClient.TransactionByHash(context.TODO(), txHash)
		printInfo(fmt.Sprintf("Transaction nonce: %d Transaction pending: %v", tx.Nonce(), isPending))
		// check transaction receipt
		receiptChan := make(chan *types.Receipt)
		client.CheckTransaction(context.TODO(), receiptChan, txHash, 1)
		receipt := <-receiptChan
		printInfo(fmt.Sprintf("Transaction status: %v", receipt.Status))
	}
}
