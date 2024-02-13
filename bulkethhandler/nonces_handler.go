package bulkethhandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Roman2K/txstreet-bulk-eth-api/limits"
	"github.com/ethereum/go-ethereum/common"
)

func newNoncesHandler(ethClient EthClient, limiter limits.Limiter) http.Handler {
	return collectAllHandler[noncesRequest, common.Address, nonceResult]{
		limiter:        limiter,
		collectAdapter: noncesAdapter{ethClient},
	}
}

type noncesRequest struct {
	Accounts []common.Address `json:"accounts"`
}

type nonceResult struct {
	Account common.Address `json:"account"`
	Count   uint64         `json:"count"`
}

type noncesAdapter struct {
	ethClient EthClient
}

func (noncesAdapter) getInputs(req noncesRequest) []common.Address {
	return req.Accounts
}

func (a noncesAdapter) collectFunc(
	ctx context.Context, account common.Address,
) (
	result nonceResult, err error,
) {
	result.Account = account
	result.Count, err = a.ethClient.NonceAt(ctx, account, nil)

	if err != nil {
		err = fmt.Errorf("Failed to get nonce of account %s: %w", account, err)
	}

	return
}
