package blocklytics

import (
	"github.com/sc0Vu/graphql"
)

// Block represents graphql model of block
type Block struct {
	ID               graphql.String
	ParentHash       graphql.String
	UnclesHash       graphql.String
	StateRoot        graphql.String
	ReceiptsRoot     graphql.String
	TransactionsRoot graphql.String
	Number           graphql.BigInt
	Timestamp        graphql.BigInt
	Author           graphql.String
	Difficulty       graphql.BigInt
	TotalDifficulty  graphql.BigInt
	GasUsed          graphql.BigInt
	GasLimit         graphql.BigInt
	Size             graphql.BigInt
}
