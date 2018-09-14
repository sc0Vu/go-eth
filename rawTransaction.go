package main

import (
	"context"
	"fmt"
	"math/big"
	"go-eth/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	if client, err := eth.Connect("http://localhost:8545"); err != nil {
		fmt.Errorf(err.Error())
		return
	} else {
		// generate private key
		// privKey, err := crypto.GenerateKey()
		// sha3 helloeth
		if privKey, err := crypto.HexToECDSA("14c8e3bfacd31c7dddba84c0ba51a2d45fb1dd299bcb9772487232a7c3d18012"); err != nil {
			fmt.Errorf(err.Error())
			return
		} else {
		    from := crypto.PubkeyToAddress(privKey.PublicKey)
			to := common.HexToAddress("30b82c8694b59695d78f33a7ba1c2a55dfa618d5")
			fmt.Println(from.String())
			if nonce, err := client.EthClient.NonceAt(context.TODO(), from, nil); err != nil {
				fmt.Errorf(err.Error())
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
					fmt.Errorf(err.Error())
				} else {
					receiptChan := make(chan *types.Receipt)
					fmt.Printf("Transaction hash: %s\n", txHash.String())
					_, isPending, _ := client.EthClient.TransactionByHash(context.TODO(), txHash)
					fmt.Printf("Transaction pending: %v\n", isPending)
					// check transaction receipt
					client.CheckTransaction(context.TODO(), receiptChan, txHash, 1)
					receipt := <-receiptChan
					fmt.Printf("Transaction status: %v\n", receipt.Status)
				}
			}
		}
	}
}