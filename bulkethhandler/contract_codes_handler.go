package bulkethhandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/limits"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func newContractCodesHandler(ethClient EthClient, limiter limits.Limiter) http.Handler {
	return collectAllHandler[contractCodesRequest, common.Address, contractCodeResult]{
		limiter:        limiter,
		collectAdapter: contractCodesAdapter{ethClient},
	}
}

type contractCodesRequest struct {
	Contracts []common.Address `json:"contracts"`
}

type contractCodeResult struct {
	Contract common.Address `json:"contract"`
	Code     string         `json:"code"`
}

type contractCodesAdapter struct {
	ethClient EthClient
}

func (contractCodesAdapter) getInputs(req contractCodesRequest) []common.Address {
	return req.Contracts
}

func (a contractCodesAdapter) collectFunc(
	ctx context.Context, account common.Address,
) (
	result contractCodeResult, err error,
) {
	result.Contract = account

	var code []byte
	code, err = a.ethClient.CodeAt(ctx, account, nil)
	result.Code = hexutil.Encode(code)

	if err != nil {
		err = fmt.Errorf(
			"Failed to get contract code of account %s: %w", account, err,
		)
	}

	return
}
