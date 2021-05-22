package main

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/sc0Vu/subgraph/tokenlon"
	"github.com/sc0Vu/subgraph/uniswap"
)

func newTokenlonClient() (cli tokenlon.TokenlonClient) {
	cli = tokenlon.NewTokenlonClient("")
	return
}

func newUniswapV2Client() (cli uniswap.UniswapV2Client) {
	cli = uniswap.NewUniswapV2Client("")
	return
}

func newCtx() (ctx context.Context) {
	ctx = context.Background()
	return
}

func main() {
	etherscanToken := os.Getenv("ETHERSCAN_TOKEN")
	ecli := tokenlon.NewEtherscan(etherscanToken)
	uc := newUniswapV2Client()
	c := newTokenlonClient()
	ctx := newCtx()
	swappeds, err := c.Swappeds(ctx, 1000, 0, "feeFactor desc")
	if err != nil {
		fmt.Println(err)
		return
	}
	var totalFeeEthCollected float64
	// sleep 1 millisecond every five requests
	for i, swapped := range swappeds {
		rf := new(big.Int)
		gas := new(big.Int)
		nf := new(big.Int)
		ra := swapped.ReceivedAmount.Int
		sa := swapped.SettleAmount.Int
		// calculate raw fee
		rf.Sub(ra, sa)
		// ethereum miner fee
		txReceipt, err := ecli.GetTransactionReceipt(newCtx(), string(swapped.TransactionHash))
		if err != nil {
			fmt.Printf("Failed to get transaction %s: %+v\n", string(swapped.TransactionHash), err)
			continue
		}
		gasUsed := new(big.Int)
		gasUsed.SetString(txReceipt.GasUsed[2:], 16)
		gas.Mul(gasUsed, swapped.GasPrice.Int)
		// fromEth := swapped.TakerAssetAddr == "0x0000000000000000000000000000000000000000"
		toEth := swapped.MakerAssetAddr == "0x0000000000000000000000000000000000000000"
		fmt.Printf("[%v] Transaction: %s\nGas fee (eth): %f\n", toEth, swapped.ID, float64(gas.Int64())/math.Pow10(18))
		var rawFeeEth float64
		var feeEth float64
		if toEth {
			// collected eth as fee
			rawFeeEth = float64(rf.Int64()) / math.Pow10(18)
			nf.Sub(rf, gas)
			feeEth = float64(nf.Int64()) / math.Pow10(18)
		} else {
			bn := swapped.BlockNumber.Int
			// fetch token price from uniswap
			mtoken, err := uc.TokensWithBN(newCtx(), string(swapped.MakerAssetAddr), int(bn.Int64()))
			if err != nil {
				fmt.Printf("Failed to fetch the token %s at %d\n", swapped.TakerAssetAddr, bn.Int64())
				continue
			}
			// ethPrice, err := uc.BundlesWithBN(ccc, 1, int(bn.Int64()))
			// if err != nil {
			// 	fmt.Printf("Failed to fetch the bundle at %d\n", bn.Int64())
			// 	continue
			// }
			tpEth, err := strconv.ParseFloat(string(mtoken.DerivedETH), 64)
			if err != nil {
				fmt.Printf("Failed to parse the token price %s at %d: %+v\n", swapped.MakerAssetAddr, bn.Int64(), err)
				continue
			}
			fmt.Printf("Token price at %d: %f\n", bn.Int64(), tpEth)
			decimal, err := strconv.ParseInt(string(mtoken.Decimals), 10, 64)
			if err != nil {
				fmt.Printf("Failed to parse the token decimal %s at %d\n", swapped.TakerAssetAddr, bn.Int64())
				continue
			}
			// fmt.Printf("Token price to eth: %f %f %d \n", float64(rf.Int64()), tpEth, decimal)
			rawFeeEth = float64(rf.Int64()) / math.Pow10(int(decimal)) * tpEth
			feeEth = rawFeeEth - (float64(gas.Int64()) / math.Pow10(18))
		}
		fmt.Printf("Raw fee collected (eth): %f\n", rawFeeEth)
		fmt.Printf("Fee collected (eth): %f\n", feeEth)
		totalFeeEthCollected += feeEth
		if i > 0 && i%5 == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}
	fmt.Printf("Total fee collected (eth): %f\n", totalFeeEthCollected)
}
