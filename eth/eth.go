package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	rpcClient *rpc.Client
	EthClient *ethclient.Client
}

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("0x%x", n)
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
func (ec *Client) GetBlockNumber(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := ec.rpcClient.CallContext(ctx, &result, "eth_blockNumber")
	return (*big.Int)(&result), err
}

// Message is a fully derived transaction and implements core.Message
type Message struct {
	To       *common.Address `json:"to"`
	From     common.Address  `json:"from"`
	Value    string          `json:"value"`
	GasLimit string          `json:"gas"`
	GasPrice string          `json:"gasPrice"`
	Data     []byte          `json:"data"`
}

// NewMessage returns the message.
func NewMessage(from common.Address, to *common.Address, value *big.Int, gasLimit *big.Int, gasPrice *big.Int, data []byte) Message {
	return Message{
		From:     from,
		To:       to,
		Value:    toHexInt(value),
		GasLimit: toHexInt(gasLimit),
		GasPrice: toHexInt(gasPrice),
		Data:     data,
	}
}

// SendTransaction injects a transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ec *Client) SendTransaction(ctx context.Context, tx *Message) (common.Hash, error) {
	var txHash common.Hash
	err := ec.rpcClient.CallContext(ctx, &txHash, "eth_sendTransaction", tx)
	return txHash, err
}
