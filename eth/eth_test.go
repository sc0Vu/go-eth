package eth

import (
	"testing"
)

func TestConnect(t *testing.T) {
	_, err := Connect("http://localhost:8545")

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetBlockNumber(t *testing.T) {
	_, err := GetBlockNumber()

	if err != nil {
		t.Errorf(err.Error())
	}
}
