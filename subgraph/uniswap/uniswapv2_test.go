package uniswap

import (
	"context"
	"testing"
)

const (
	targetV2BN = 11111111
	bundleID   = 1
	// test pair id for WBTC-ETH
	pairID = "0xbb2b8038a1640196fbe3e38816f3e67cba72d940"
	// test token id for uniswap
	tokenID = "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"
)

func newV2Client() (uniCli UniswapV2Client) {
	uniCli = NewUniswapV2Client("")
	return
}

func newCtx() (ctx context.Context) {
	ctx = context.Background()
	return
}

// TestBundles test subgraph bundles api
func TestBundles(t *testing.T) {
	c := newV2Client()
	ctx := newCtx()
	ethPriceNow, err := c.Bundles(ctx, bundleID)
	if err != nil {
		t.Fatal(err)
	}
	if ethPriceNow < 0 {
		t.Fatalf("Eth price should not be zero")
	}
	ethPriceOld, err := c.BundlesWithBN(ctx, bundleID, targetV2BN)
	if err != nil {
		t.Fatal(err)
	}
	if ethPriceOld < 0 {
		t.Fatalf("Eth price should not be zero")
	}
	if ethPriceNow == ethPriceOld {
		t.Fatalf("Eth price old should not be the same with eth price now")
	}
}

// TestTokens test subgraph tokens api
func TestTokens(t *testing.T) {
	c := newV2Client()
	ctx := newCtx()
	tokenNow, err := c.Tokens(ctx, tokenID)
	if err != nil {
		t.Fatal(err)
	}
	if tokenNow.ID != tokenID {
		t.Fatalf("Token ID not equal")
	}
	tokenOld, err := c.TokensWithBN(ctx, tokenID, targetV2BN)
	if err != nil {
		t.Fatal(err)
	}
	if tokenOld.ID != tokenID {
		t.Fatalf("Token ID not equal")
	}
	if tokenNow.Name != tokenOld.Name {
		t.Fatalf("Tokens name not equal")
	}
	if tokenNow.Symbol != tokenOld.Symbol {
		t.Fatalf("Tokens symbol not equal")
	}
}

// TestPairs test subgraph pairs api
func TestPairs(t *testing.T) {
	c := newV2Client()
	ctx := newCtx()
	pairNow, err := c.Pairs(ctx, pairID)
	if err != nil {
		t.Fatal(err)
	}
	if pairNow.ID != pairID {
		t.Fatalf("Pair ID not equal")
	}
	pairOld, err := c.PairsWithBN(ctx, pairID, targetV2BN)
	if err != nil {
		t.Fatal(err)
	}
	if pairOld.ID != pairID {
		t.Fatalf("Pair ID not equal")
	}
	if pairNow.Token0.Name != pairOld.Token0.Name {
		t.Fatalf("Pairs token0 name not equal")
	}
	if pairNow.Token0.Symbol != pairOld.Token0.Symbol {
		t.Fatalf("Pairs token0 name not equal")
	}
	if pairNow.Token1.Name != pairOld.Token1.Name {
		t.Fatalf("Pairs token1 name not equal")
	}
	if pairNow.Token1.Symbol != pairOld.Token1.Symbol {
		t.Fatalf("Pairs token1 name not equal")
	}
}
