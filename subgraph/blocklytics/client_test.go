package blocklytics

import (
	"context"
	"testing"
	"time"
)

const (
	targetBN = 11111111
)

func newClient() (cli BlocklyticsClient) {
	cli = NewBlocklyticsClient("")
	return
}

func newCtx() (ctx context.Context) {
	ctx = context.Background()
	return
}

// TestBlocks test subgraph blocks api
func TestBlocks(t *testing.T) {
	c := newClient()
	ctx := newCtx()
	blocks, err := c.Blocks(ctx, 10, 0, targetBN, "number desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 10 {
		t.Fatalf("Should fetch 10 blocks")
	}
	for _, block := range blocks {
		number := int(block.Number.Int.Int64())
		if number < targetBN {
			t.Fatalf("Should fetch blocks that block number greater than %d", targetBN)
		}
	}
}

// TestBlocksByTimestamp test subgraph blocks api
func TestBlocksByTimestamp(t *testing.T) {
	c := newClient()
	ctx := newCtx()
	oneDay := time.Now().Add(-86400 * time.Second)
	targetTP := int(oneDay.Unix())
	blocks, err := c.BlocksByTimestamp(ctx, 10, 0, targetTP, "timestamp asc")
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 10 {
		t.Fatalf("Should fetch 10 blocks")
	}
	for _, block := range blocks {
		timestamp := int(block.Timestamp.Int.Int64())
		if int(timestamp) < targetTP {
			t.Fatalf("Should fetch blocks that block number greater than %d", targetTP)
		}
	}
}
