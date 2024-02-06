package bulkethhandler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/limits"
)

type handler struct {
	ctx       context.Context
	ethClient EthClient
	limiter   limits.Limiter
}

func NewHandler(ctx context.Context, ethClient EthClient, limiter limits.Limiter) http.Handler {
	handler := handler{
		ctx:       ctx,
		ethClient: ethClient,
		limiter:   limiter,
	}

	mux := http.NewServeMux()

	//
	//
	// /transaction-receipts

	mux.HandleFunc("/nonces", handler.handleNonces)
	mux.HandleFunc("/contract-codes", handler.handleContractCodes)

	return mux
}

func sendError(w http.ResponseWriter, err error) {
	slog.Error("Error while serving request", "error", err)

	http.Error(w, "Internal error", 500)
}
