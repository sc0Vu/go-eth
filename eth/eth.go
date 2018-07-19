package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	rpcClient *rpc.Client
	EthClient *ethclient.Client
}

// Connect creates a client that uses the given host.
func Connect(host string) (*Client, error) {
	rpcClient, err := rpc.Dial(host)

	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)

	return &Client{rpcClient, ethClient}, nil
}

// GetBlockNumber returns the block number.
func (ec *Client) GetBlockNumber() (*big.Int, error) {
	var result hexutil.Big
	err := ec.rpcClient.CallContext(context.TODO(), &result, "eth_blockNumber")
	return (*big.Int)(&result), err
}
