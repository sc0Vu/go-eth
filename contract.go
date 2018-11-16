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

func main() {
	if client, err := eth.Connect("ws://localhost:8546"); err != nil {
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
			fmt.Println(from.String())
			if nonce, err := client.EthClient.NonceAt(context.TODO(), from, nil); err != nil {
				fmt.Errorf(err.Error())
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
				// 	fmt.Errorf(err.Error())
				// } else {
				// 	fmt.Println(contractAddress, deployTx.Hash().String(), contract)
				// }
				// EIP155 signer
				// signer := types.NewEIP155Signer(big.NewInt(4))
				signer := types.HomesteadSigner{}
				signedTx, _ := types.SignTx(tx, signer, privKey)
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
					fmt.Printf("Contract address: %s\n", receipt.ContractAddress.Hex())

					// get balance
					contract, _ := Contract.NewContract(receipt.ContractAddress, client.EthClient)
					balance, _ := contract.BalanceOf(nil, from)
					fmt.Printf("balance of %s is %d\n", from.String(), balance)

					// watch transfer
					ch := make(chan *Contract.ContractTransfer)
					sub, _ := contract.WatchTransfer(nil, ch, nil, nil)

					go func() {
						for {
							select {
							case err := <-sub.Err():
								fmt.Errorf(err.Error())
								os.Exit(1)
							case log := <-ch:
								fmt.Printf("[Transfer event]\n")
								fmt.Printf("From: %s\n", log.From.String())
								fmt.Printf("To: %s\n", log.To.String())
								fmt.Printf("Value: %d\n", log.Value)
								os.Exit(0)
							}
						}
					}()
					fmt.Println("Event watching...")

					// transfer token
					nonce++
					non := big.NewInt(int64(nonce))
					txOpts := makeTxOpts(from, non, amount, gasPrice, gasLimit, privKey, 0)
					toAddress := common.HexToAddress("0x9b23a6a9a60b3846f86ebc451d11bef20ed07930")
					bal := big.NewInt(10000)

					if transferTx, err := contract.Transfer(txOpts, toAddress, bal); err != nil {
						fmt.Errorf(err.Error())
					} else {
						txHash := transferTx.Hash()
						receiptChan := make(chan *types.Receipt)
						fmt.Printf("Transaction hash: %s\n", txHash.String())
						_, isPending, _ := client.EthClient.TransactionByHash(context.TODO(), txHash)
						fmt.Printf("Transaction pending: %v\n", isPending)
						// check transaction receipt
						client.CheckTransaction(context.TODO(), receiptChan, txHash, 1)
						receipt := <-receiptChan
						fmt.Printf("Transaction status: %v\n", receipt.Status)
					}
					<-ch
				}
			}
		}
	}
}
