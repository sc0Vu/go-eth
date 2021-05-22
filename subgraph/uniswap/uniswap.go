package uniswap

import (
	"github.com/sc0Vu/graphql"
)

// UniswapClient represents uniswapv2 subgraph client
type UniswapClient struct {
	c     *graphql.Client
	token string
}

// NewUniswapClient returns uniswap subgraph client
func NewUniswapClient(token string) (uniCli UniswapClient) {
	uniCli.c = graphql.NewClient(v1Endpoint, nil)
	uniCli.token = token
	return
}
