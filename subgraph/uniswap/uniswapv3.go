package uniswap

import (
	"context"
	"strconv"

	"github.com/sc0Vu/graphql"
)

// TODO: remove reflect usage?
// TODO: big decimal type in case of precision?

// UniswapV3Client represents uniswapv3 subgraph client
type UniswapV3Client struct {
	UniswapClient
}

// NewUniswapV3Client returns uniswapv2 subgraph client
func NewUniswapV3Client(token string) (uniCli UniswapV3Client) {
	uniCli.c = graphql.NewClient(v3Endpoint, nil)
	uniCli.token = token
	return
}

// BundlesWithBN returns the price of eth in given block number
func (uniCli *UniswapV3Client) BundlesWithBN(ctx context.Context, id, bn int) (ethPrice float64, err error) {
	var query struct {
		Bundle V3Bundle `graphql:"bundle(id: $id, block:{number: $bn})"`
	}
	variables := map[string]interface{}{
		"id": graphql.Int(id),
		"bn": graphql.Int(bn),
	}
	err = uniCli.c.Query(ctx, &query, variables)
	if err != nil {
		return
	}
	if ethPrice, err = strconv.ParseFloat(string(query.Bundle.EthPriceUSD), 64); err != nil {
		return
	}
	return
}

// Bundles returns the current price of eth
func (uniCli *UniswapV3Client) Bundles(ctx context.Context, id int) (ethPrice float64, err error) {
	var query struct {
		Bundle V3Bundle `graphql:"bundle(id: $id)"`
	}
	variables := map[string]interface{}{
		"id": graphql.Int(id),
	}
	err = uniCli.c.Query(ctx, &query, variables)
	if err != nil {
		return
	}
	if ethPrice, err = strconv.ParseFloat(string(query.Bundle.EthPriceUSD), 64); err != nil {
		return
	}
	return
}

// TokensWithBN returns the token in uniswap of given address and block number
func (uniCli *UniswapV3Client) TokensWithBN(ctx context.Context, id string, bn int) (token V3Token, err error) {
	var query struct {
		Token V3Token `graphql:"token(id: $id, where:{id: $tokenID}, block:{number: $bn})"`
	}
	variables := map[string]interface{}{
		"id":      graphql.String(id),
		"tokenID": graphql.String(id),
		"bn":      graphql.Int(bn),
	}
	err = uniCli.c.Query(ctx, &query, variables)
	if err != nil {
		return
	}
	token = query.Token
	return
}

// Tokens returns the token in uniswap of given address
func (uniCli *UniswapV3Client) Tokens(ctx context.Context, id string) (token V3Token, err error) {
	var query struct {
		Token V3Token `graphql:"token(id: $id, where:{id: $id})"`
	}
	variables := map[string]interface{}{
		"id": graphql.String(id),
	}
	err = uniCli.c.Query(ctx, &query, variables)
	if err != nil {
		return
	}
	token = query.Token
	return
}
