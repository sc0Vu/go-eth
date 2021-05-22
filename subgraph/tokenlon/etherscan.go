package tokenlon

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const etherscanEndpoint = "https://api.etherscan.io/api"

var ErrWrongResult = fmt.Errorf("wrong result")

type Response struct {
	Message string
	Status  string
	Result  TransactionReceipt
}

type Etherscan struct {
	token    string
	endpoint string
}

type TransactionReceipt struct {
	GasUsed string
}

func NewEtherscan(token string) (esc Etherscan) {
	esc.token = token
	esc.endpoint = etherscanEndpoint
	return esc
}

// GetTransactionReceipt returns transaction receipt of given transaction hash
// https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=0x0566aeb88af83608b6918a328ddb36a0b0207ad9b72d89f453df3d780ae7a231&apikey=T67NS47W56F7MKDR38UHAZTIJ5ZHEZ7FYE
func (esc *Etherscan) GetTransactionReceipt(ctx context.Context, txHash string) (TransactionReceipt, error) {
	url := fmt.Sprintf("%s?module=%s&action=%s&txhash=%s&apikey=%s", esc.endpoint, "proxy", "eth_getTransactionReceipt", txHash, esc.token)
	c := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return TransactionReceipt{}, err
	}

	res, err := c.Do(req)
	if err != nil {
		return TransactionReceipt{}, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TransactionReceipt{}, err
	}

	var apiRes Response
	err = json.Unmarshal(body, &apiRes)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return apiRes.Result, nil
}
