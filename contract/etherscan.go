package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const etherscanEndpoint = "https://api.etherscan.io/api"

var ErrWrongResult = fmt.Errorf("wrong result")

type Response struct {
	Message string
	Status  string
	Result  interface{}
}

type Etherscan struct {
	token    string
	endpoint string
}

func NewEtherscan(token string) (esc Etherscan) {
	esc.token = token
	esc.endpoint = etherscanEndpoint
	return esc
}

func (esc *Etherscan) GetVerifiedABI(ctx context.Context, contractAddr common.Address) (string, error) {
	url := fmt.Sprintf("%s?module=%s&action=%s&address=%s&apikey=%s", esc.endpoint, "contract", "getabi", contractAddr.String(), esc.token)
	c := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := c.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var apiRes Response
	err = json.Unmarshal(body, &apiRes)
	if err != nil {
		return "", err
	}
	if apiRes.Status != "1" || apiRes.Message != "OK" {
		return "", fmt.Errorf("fetch error: %s", apiRes.Message)
	}

	strABI, ok := apiRes.Result.(string)
	if !ok {
		return "", ErrWrongResult
	}
	return strABI, nil
}
