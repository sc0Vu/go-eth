package eth

import (
	"context"
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
	Nonce    uint64          `json:"nonce"`
	Value    *big.Int        `json:"value"`
	GasLimit uint64          `json:"gas"`
	GasPrice *big.Int        `json:"gasPrice"`
	Data     []byte          `json:"data"`
}

// NewMessage returns the message.
func NewMessage(from common.Address, to *common.Address, nonce uint64, value *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) Message {
	return Message{
		From:     from,
		To:       to,
		Nonce:    nonce,
		Value:    value,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}
}

// SendTransaction injects a transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ec *Client) SendTransaction(ctx context.Context, tx *Message) error {
	err := ec.rpcClient.CallContext(ctx, nil, "eth_sendTransaction", tx)
	return err
}
