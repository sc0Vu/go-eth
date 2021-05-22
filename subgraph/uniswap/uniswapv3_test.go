package uniswap

import (
	"testing"
)

const (
	targetV3BN = 12411111
)

func newV3Client() (uniCli UniswapV3Client) {
	uniCli = NewUniswapV3Client("")
	return
}

// TestV3Bundles test subgraph bundles api
func TestV3Bundles(t *testing.T) {
	c := newV3Client()
	ctx := newCtx()
	ethPriceNow, err := c.Bundles(ctx, bundleID)
	if err != nil {
		t.Fatal(err)
	}
	if ethPriceNow < 0 {
		t.Fatalf("Eth price should not be zero")
	}
	ethPriceOld, err := c.BundlesWithBN(ctx, bundleID, targetV3BN)
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

// TestV3Tokens test subgraph tokens api
func TestV3Tokens(t *testing.T) {
	c := newV3Client()
	ctx := newCtx()
	tokenNow, err := c.Tokens(ctx, tokenID)
	if err != nil {
		t.Fatal(err)
	}
	if tokenNow.ID != tokenID {
		t.Fatalf("Token ID not equal")
	}
	tokenOld, err := c.TokensWithBN(ctx, tokenID, targetV3BN)
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
