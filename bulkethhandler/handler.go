package bulkethhandler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/limits"
)

func NewHandler(ctx context.Context, ethClient EthClient, limiter limits.Limiter) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/nonces", newNoncesHandler(ethClient, limiter))
	mux.Handle("/contract-codes", newContractCodesHandler(ethClient, limiter))
	mux.Handle("/transaction-receipts", newTxReceiptsHandler(ethClient, limiter))

	return mux
}

func sendError(w http.ResponseWriter, err error) {
	slog.Error("Error while serving request", "error", err)

	http.Error(w, "Internal error", 500)
}
