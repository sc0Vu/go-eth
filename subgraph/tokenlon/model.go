package tokenlon

import (
	"github.com/sc0Vu/graphql"
)

// Swapped represents graphql model of swapped
type Swapped struct {
	ID               graphql.String
	Source           graphql.String
	TransactionHash  graphql.String
	UserAddr         graphql.String
	MakerAssetAddr   graphql.String
	TakerAssetAddr   graphql.String
	TakerAssetAmount graphql.BigInt
	MakerAssetAmount graphql.BigInt
	ReceivedAmount   graphql.BigInt
	SettleAmount     graphql.BigInt
	// GasUsed          graphql.BigInt
	GasPrice    graphql.BigInt
	BlockNumber graphql.BigInt
	FeeFactor   graphql.Int
}
