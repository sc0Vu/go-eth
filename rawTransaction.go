package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"go-eth/eth"
	"math/big"
	"os"
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
		// generate private key
		// privKey, err := crypto.GenerateKey()
		// sha3 helloeth
		if privKey, err := crypto.HexToECDSA("14c8e3bfacd31c7dddba84c0ba51a2d45fb1dd299bcb9772487232a7c3d18012"); err != nil {
			printError(err)
			return
		} else {
			from := crypto.PubkeyToAddress(privKey.PublicKey)
			to := common.HexToAddress("30b82c8694b59695d78f33a7ba1c2a55dfa618d5")
			printInfo(from.String())
			if nonce, err := client.EthClient.NonceAt(context.TODO(), from, nil); err != nil {
				printError(err)
				return
			} else {
				amount := big.NewInt(1)
				gasLimit := uint64(90000)
				gasPrice := big.NewInt(0)
				data := []byte{}
				tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
				// EIP155 signer
				// signer := types.NewEIP155Signer(big.NewInt(4))
				signer := types.HomesteadSigner{}
				signedTx, _ := types.SignTx(tx, signer, privKey)
				// client.EthClient.SendTransaction(context.TODO(), signedTx)
				if txHash, err := client.SendRawTransaction(context.TODO(), signedTx); err != nil {
					printError(err)
				} else {
					receiptChan := make(chan *types.Receipt)
					printInfo(fmt.Sprintf("Transaction hash: %s", txHash.String()))
					_, isPending, _ := client.EthClient.TransactionByHash(context.TODO(), txHash)
					printInfo(fmt.Sprintf("Transaction pending: %v", isPending))
					// check transaction receipt
					client.CheckTransaction(context.TODO(), receiptChan, txHash, 1)
					receipt := <-receiptChan
					printInfo(fmt.Sprintf("Transaction status: %v", receipt.Status))
				}
			}
		}
	}
}
