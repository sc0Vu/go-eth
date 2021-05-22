package blocklytics

import (
	"context"
	"fmt"
	"github.com/sc0Vu/graphql"
	"strings"
)

const endpoint = "https://api.thegraph.com/subgraphs/name/blocklytics/ethereum-blocks"

type BlocklyticsClient struct {
	c     *graphql.Client
	token string
}

// TODO: fix id: 1
// TODO: query maker?

// NewBlocklyticsClient returns Blocklytics
func NewBlocklyticsClient(token string) (bl BlocklyticsClient) {
	bl.c = graphql.NewClient(endpoint, nil)
	bl.token = token
	return
}

// Blocks query of blocks by given query, the order should looks like [column asc|desc]
func (bl *BlocklyticsClient) Blocks(ctx context.Context, count, skip, bn int, order string) (blocks []Block, err error) {
	pOrder := strings.Split(order, " ")
	if len(pOrder) != 2 {
		err = fmt.Errorf("the order should looks like [column asc|desc]")
		return
	}
	if pOrder[1] != "asc" && pOrder[1] != "desc" {
		err = fmt.Errorf("the order direction should be one of asc or desc")
		return
	}
	var query struct {
		Blocks []Block `graphql:"blocks(id: 1, first: $count, skip: $skip, orderBy: $orderBy, orderDirection: $orderDir, where: {number_gt: $bn})"`
	}
	variables := map[string]interface{}{
		"count":    graphql.Int(count),
		"skip":     graphql.Int(skip),
		"bn":       graphql.Int(bn),
		"orderBy":  graphql.String(pOrder[0]),
		"orderDir": graphql.String(pOrder[1]),
	}
	err = bl.c.Query(ctx, &query, variables)
	if err != nil {
		return
	}
	blocks = query.Blocks
	return
}

// BlocksByTimestamp query of blocks by given query, the order should looks like [column asc|desc]
func (bl *BlocklyticsClient) BlocksByTimestamp(ctx context.Context, count, skip, timestamp int, order string) (blocks []Block, err error) {
	pOrder := strings.Split(order, " ")
	if len(pOrder) != 2 {
		err = fmt.Errorf("the order should looks like [column asc|desc]")
		return
	}
	if pOrder[1] != "asc" && pOrder[1] != "desc" {
		err = fmt.Errorf("the order direction should be one of asc or desc")
		return
	}
	var query struct {
		Blocks []Block `graphql:"blocks(id: 1, first: $count, skip: $skip, orderBy: $orderBy, orderDirection: $orderDir, where: {timestamp_gt: $timestampFrom, timestamp_lt: $timestampTo})"`
	}
	variables := map[string]interface{}{
		"count":         graphql.Int(count),
		"skip":          graphql.Int(skip),
		"timestampFrom": graphql.Int(timestamp),
		"timestampTo":   graphql.Int(timestamp + 600),
		"orderBy":       graphql.String(pOrder[0]),
		"orderDir":      graphql.String(pOrder[1]),
	}
	err = bl.c.Query(ctx, &query, variables)
	if err != nil {
		return
	}
	blocks = query.Blocks
	return
}
