package uniswap

import (
	"github.com/sc0Vu/graphql"
)

// Bundle represent graphql model of Bundle in uniswap v2
type Bundle struct {
	ID       graphql.ID
	EthPrice graphql.String
}

// Token represent graph model of token in uniswap v2
type Token struct {
	ID                 graphql.ID
	Symbol             graphql.String
	Name               graphql.String
	Decimals           graphql.String
	TotalSupply        graphql.String
	TradeVolume        graphql.String
	TradeVolumeUSD     graphql.String `graphql:"tradeVolumeUSD"`
	UntrackedVolumeUSD graphql.String `graphql:"untrackedVolumeUSD"`
	TXCount            graphql.String
	TotalLiquidity     graphql.String
	DerivedETH         graphql.String `graphql:"derivedETH"`
}

// Pair represent graphql model of Pair in uniswap v2
type Pair struct {
	ID                     graphql.ID
	Token0                 Token
	Token1                 Token
	TrackedReserveETH      graphql.String `graphql:"trackedReserveETH"`
	VolumeUSD              graphql.String `graphql:"volumeUSD"`
	UntrackedVolumeUSD     graphql.String `graphql:"untrackedVolumeUSD"`
	TXCount                graphql.String `graphql:"txCount"`
	CreatedAtTimestamp     graphql.String `graphql:"createdAtTimestamp"`
	CreatedAtBlockNumber   graphql.String `graphql:"createdAtBlockNumber"`
	LiquidityProviderCount graphql.String
	Reserve0               graphql.String
	Reserve1               graphql.String
}

// V3Bundle represent graphql model of Bundle in uniswap v3
type V3Bundle struct {
	ID          graphql.ID
	EthPriceUSD graphql.String `graphql:"ethPriceUSD"`
}

// V3Token represent graph model of token in uniswap v3
type V3Token struct {
	ID                 graphql.ID
	Symbol             graphql.String
	Name               graphql.String
	Decimals           graphql.String
	TotalSupply        graphql.String
	Volume             graphql.String
	VolumeUSD          graphql.String `graphql:"volumeUSD"`
	UntrackedVolumeUSD graphql.String `graphql:"untrackedVolumeUSD"`
	TXCount            graphql.String
	TotalValueLocked   graphql.String
	DerivedETH         graphql.String `graphql:"derivedETH"`
}
