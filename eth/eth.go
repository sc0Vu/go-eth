package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var ctx *rpc.Client

// Connect creates a client that uses the given host.
func Connect(host string) (*ethclient.Client, error) {
	context, err := rpc.Dial(host)

	if err != nil {
		return nil, err
	}
	ctx = context
	conn := ethclient.NewClient(context)

	return conn, nil
}

// GetBlockNumber returns the block number.
func GetBlockNumber() (*big.Int, error) {
	var result hexutil.Big
	err := ctx.CallContext(context.TODO(), &result, "eth_blockNumber")
	return (*big.Int)(&result), err
}
