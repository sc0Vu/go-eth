package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sc0Vu/subgraph/blocklytics"
	"github.com/sc0Vu/subgraph/uniswap"
)

const dayseconds = 86400

func fetchBlockNumber(token string) (int64, error) {
	now := time.Now().Truncate(time.Minute)
	oneDay := now.Add(-1 * dayseconds * time.Second)
	cli := blocklytics.NewBlocklyticsClient(token)
	blocks, err := cli.BlocksByTimestamp(context.TODO(), 1, 0, int(oneDay.Unix()), "timestamp asc")
	if err != nil {
		return 0, err
	}
	if len(blocks) <= 0 {
		return 0, fmt.Errorf("couldn't find any block")
	}
	return blocks[0].Number.Int.Int64(), nil
}

func main() {
	token := os.Getenv("SUBGRAPH_TOKEN")
	pairAddress := os.Getenv("UNISWAP_ADDRESS")
	bnI64, err := fetchBlockNumber(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	bn := int(bnI64)
	cli := uniswap.NewUniswapV2Client(token)
	ethPrice, err := cli.BundlesWithBN(context.TODO(), 1, bn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Ethereum price in %d: %f\n", bn, ethPrice)
	ethPriceNow, err := cli.Bundles(context.TODO(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Current ethereum price: %f\n", ethPriceNow)

	pair, err := cli.PairsWithBN(context.TODO(), pairAddress, bn)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("Pair information in %d: %+v\n", bn, pair)
	pairNow, err := cli.Pairs(context.TODO(), pairAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("Current pair information: %+v\n", pairNow)
	reserveETHNow, err := strconv.ParseFloat(string(pairNow.TrackedReserveETH), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	liquidityNow := reserveETHNow * ethPriceNow
	fmt.Printf("Current liquidity: %f\n", liquidityNow)
	volumeNow, err := strconv.ParseFloat(string(pairNow.VolumeUSD), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	volume, err := strconv.ParseFloat(string(pair.VolumeUSD), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	oneDayVolume := volumeNow - volume
	fmt.Printf("One day volume: %f\n", oneDayVolume)
	fee := oneDayVolume * 0.003
	fmt.Printf("Fee: %f\n", fee)
	reserve0, err := strconv.ParseFloat(string(pairNow.Reserve0), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	reserve1, err := strconv.ParseFloat(string(pairNow.Reserve1), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	token0Rate := reserve1 / reserve0
	token1Rate := reserve0 / reserve1
	token0DerivedETH, err := strconv.ParseFloat(string(pairNow.Token0.DerivedETH), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	token0Price := token0DerivedETH * ethPriceNow
	fmt.Printf("1 %s = %f %s (%f USD)\n", pairNow.Token0.Symbol, token0Rate, pairNow.Token1.Symbol, token0Price)
	token1DerivedETH, err := strconv.ParseFloat(string(pairNow.Token1.DerivedETH), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	token1Price := token1DerivedETH * ethPriceNow
	fmt.Printf("1 %s = %f %s (%f USD)\n", pairNow.Token1.Symbol, token1Rate, pairNow.Token0.Symbol, token1Price)
}
