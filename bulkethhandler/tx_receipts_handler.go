package bulkethhandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/limits"
	"github.com/ethereum/go-ethereum/common"
)

func newTxReceiptsHandler(ethClient EthClient, limiter limits.Limiter) http.Handler {
	return collectAllHandler[txReceiptsRequest, common.Hash, txReceiptResult]{
		limiter:        limiter,
		collectAdapter: txReceiptsAdapter{ethClient},
	}
}

type txReceiptsRequest struct {
	Hashes []common.Hash `json:"hashes"`
}

type txReceiptResult struct {
	Hash    common.Hash `json:"hash"`
	Receipt interface{} `json:"receipt"`
}

type txReceiptsAdapter struct {
	ethClient EthClient
}

func (txReceiptsAdapter) getInputs(req txReceiptsRequest) []common.Hash {
	return req.Hashes
}

func (a txReceiptsAdapter) collectFunc(
	ctx context.Context, hash common.Hash,
) (
	result txReceiptResult, err error,
) {
	result.Hash = hash

	// We can't use `ethClient.TransactionReceipt(ctx, hash)` as types.Receipt
	// doesn't include `From` or `To`.
	err =
		a.ethClient.
			Client().
			CallContext(ctx, &result.Receipt, "eth_getTransactionReceipt", hash)

	if err != nil {
		err = fmt.Errorf("Failed to get receipt of transaction %s: %w", hash, err)
	}

	return
}
