package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	Contract "go-eth/contract"
	"go-eth/eth"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// makeTxOpts make transaction options for contrac function
func makeTxOpts(from common.Address, nonce *big.Int, value *big.Int, gasPrice *big.Int, gasLimit uint64, privKey *ecdsa.PrivateKey, chainID int64) *bind.TransactOpts {
	txOpts := &bind.TransactOpts{
		From:  from,
		Nonce: nonce,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			var txSigner types.Signer
			if chainID != 0 {
				// EIP155 signer
				txSigner = types.NewEIP155Signer(big.NewInt(chainID))
			} else {
				// default is homestead signer
				txSigner = signer
			}
			signedTx, err := types.SignTx(tx, txSigner, privKey)
			if err != nil {
				return nil, err
			}
			return signedTx, nil
		},
		Value:    value,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	}
	return txOpts
}

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
		rpc = "ws://localhost:8546"
	}
	if client, err := eth.Connect(rpc); err != nil {
		printError(err)
		return
	} else {
		// generate private key
		// privKey, err := crypto.GenerateKey()
		// sha3 helloeth
		if privKey, err := crypto.HexToECDSA("14c8e3bfacd31c7dddba84c0ba51a2d45fb1dd299bcb9772487232a7c3d18012"); err != nil {
			printError(err)
			return
		} else {
			from := crypto.PubkeyToAddress(privKey.PublicKey)
			printInfo(from.String())
			if nonce, err := client.EthClient.NonceAt(context.TODO(), from, nil); err != nil {
				printError(err)
				return
			} else {
				// deploy contract
				// you can send ether when you set payable callback function in contract
				// amount := big.NewInt(1000000000000000000)
				amount := big.NewInt(0)
				gasLimit := uint64(4600000)
				gasPrice := big.NewInt(1000000000)
				data := common.FromHex(Contract.ContractBin)
				tx := types.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)
				// another way to deploy contract
				// non := big.NewInt(int64(nonce))
				// txOpts := makeTxOpts(from, non, amount, gasPrice, gasLimit, privKey, 0)
				// if contractAddress, deployTx, contract, err := Contract.DeployContract(txOpts, client.EthClient); err != nil {
				// 	printError(err)
				// } else {
				// 	printInfo(contractAddress, deployTx.Hash().String(), contract)
				// }
				// EIP155 signer
				// signer := types.NewEIP155Signer(big.NewInt(4))
				signer := types.HomesteadSigner{}
				signedTx, _ := types.SignTx(tx, signer, privKey)
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
					printInfo(fmt.Sprintf("Contract address: %s", receipt.ContractAddress.Hex()))

					// get balance
					contract, _ := Contract.NewContract(receipt.ContractAddress, client.EthClient)
					balance, _ := contract.BalanceOf(nil, from)
					printInfo(fmt.Sprintf("balance of %s is %d", from.String(), balance))

					// watch transfer, you need to use websocket to watch event
					ch := make(chan *Contract.ContractTransfer)
					sub, _ := contract.WatchTransfer(nil, ch, nil, nil)

					go func() {
						for {
							select {
							case err := <-sub.Err():
								printError(err)
								os.Exit(1)
							case log := <-ch:
								printInfo("[Transfer event]")
								printInfo(fmt.Sprintf("From: %s", log.From.String()))
								printInfo(fmt.Sprintf("To: %s", log.To.String()))
								printInfo(fmt.Sprintf("Value: %d", log.Value))
								os.Exit(0)
							}
						}
					}()
					printInfo("Event watching...")
					// use filter transfer instead
					// filterOpts := &bind.FilterOpts{
					// 	Start: 1000000,
					// 	End:   nil,
					// }
					// if logs, err := contract.FilterTransfer(filterOpts, nil, nil); err != nil {
					// 	printError(err)
					// 	os.Exit(1)
					// } else {
					// 	log := logs.Event
					// 	printInfo("[Transfer event]\n")
					// 	printInfo(fmt.Sprintf("From: %s", log.From.String()))
					// 	printInfo(fmt.Sprintf("To: %s", log.To.String()))
					// 	printInfo(fmt.Sprintf("Value: %d", log.Value))
					// }

					// transfer token
					nonce++
					non := big.NewInt(int64(nonce))
					txOpts := makeTxOpts(from, non, amount, gasPrice, gasLimit, privKey, 0)
					toAddress := common.HexToAddress("0x9b23a6a9a60b3846f86ebc451d11bef20ed07930")
					bal := big.NewInt(10000)

					if transferTx, err := contract.Transfer(txOpts, toAddress, bal); err != nil {
						printError(err)
					} else {
						txHash := transferTx.Hash()
						receiptChan := make(chan *types.Receipt)
						printInfo(fmt.Sprintf("Transaction hash: %s", txHash.String()))
						_, isPending, _ := client.EthClient.TransactionByHash(context.TODO(), txHash)
						printInfo(fmt.Sprintf("Transaction pending: %v", isPending))
						// check transaction receipt
						client.CheckTransaction(context.TODO(), receiptChan, txHash, 1)
						receipt := <-receiptChan
						printInfo(fmt.Sprintf("Transaction status: %v", receipt.Status))
					}
					<-ch
				}
			}
		}
	}
}
