package tokenlon

import (
	"context"
	"testing"
)

func newClient() (cli TokenlonClient) {
	cli = NewTokenlonClient("")
	return
}

func newCtx() (ctx context.Context) {
	ctx = context.Background()
	return
}

// TestSwappeds test subgraph swappeds api
func TestSwappeds(t *testing.T) {
	c := newClient()
	ctx := newCtx()
	swappeds, err := c.Swappeds(ctx, 10, 0, "feeFactor desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(swappeds) != 10 {
		t.Fatalf("Should fetch 10 swappeds")
	}
}
