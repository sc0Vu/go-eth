package eth

import (
	"testing"
)

var ctx *Client

func TestConnect(t *testing.T) {
	client, err := Connect("http://localhost:8545")

	if err != nil {
		t.Errorf(err.Error())
	}
	ctx = client
}

func TestGetBlockNumber(t *testing.T) {
	_, err := ctx.GetBlockNumber()

	if err != nil {
		t.Errorf(err.Error())
	}
}
