package tokenlon

import (
	"context"
	"fmt"
	"github.com/sc0Vu/graphql"
	"strings"
)

const endpoint = "https://api.thegraph.com/subgraphs/name/consenlabs/tokenlon-v5-exchange"

type TokenlonClient struct {
	c     *graphql.Client
	token string
}

// TODO: fix id: 1
// TODO: query maker?

// NewTokenlonClient returns TokenlonClient
func NewTokenlonClient(token string) (tl TokenlonClient) {
	tl.c = graphql.NewClient(endpoint, nil)
	tl.token = token
	return
}

// Swappeds query of swappeds by given query, the order should looks like [column asc|desc]
func (tl *TokenlonClient) Swappeds(ctx context.Context, count, skip int, order string) (swappeds []Swapped, err error) {
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
		Swappeds []Swapped `graphql:"swappeds(first: $count, skip: $skip, orderBy: $orderBy, orderDirection: $orderDir)"`
	}
	variables := map[string]interface{}{
		"count":    graphql.Int(count),
		"skip":     graphql.Int(skip),
		"orderBy":  graphql.String(pOrder[0]),
		"orderDir": graphql.String(pOrder[1]),
	}
	err = tl.c.Query(ctx, &query, variables)
	if err != nil {
		return
	}
	swappeds = query.Swappeds
	return
}
