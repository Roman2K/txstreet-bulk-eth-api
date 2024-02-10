package bulkethhandler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/contextutil"
	"github.com/Roman2K/bulk-eth-api/limits"
)

type handler struct {
	parentContext context.Context
	handler       http.Handler
}

func NewHandler(ctx context.Context, ethClient EthClient, limiter limits.Limiter) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/nonces", newNoncesHandler(ethClient, limiter))
	mux.Handle("/contract-codes", newContractCodesHandler(ethClient, limiter))
	mux.Handle("/transaction-receipts", newTxReceiptsHandler(ethClient, limiter))

	return handler{
		parentContext: ctx,
		handler: requestLogHandler{
			mux,
		},
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := contextutil.Intersect(h.parentContext, r.Context())
	defer cancel()

	ctx = withRequestId(ctx, newRequestId())
	r = r.WithContext(ctx)

	h.handler.ServeHTTP(w, r)
}

func sendError(w http.ResponseWriter, r *http.Request, err error) {
	logger := requestLogger(r)
	logger.ErrorContext(r.Context(), "Error while serving request", "error", err)

	http.Error(w, "Internal error", 500)
}

func requestLogger(r *http.Request) *slog.Logger {
	logger := slog.Default()

	if requestId, ok := contextRequestId(r.Context()); ok {
		return logger.With("requestId", requestId)
	}

	return logger
}
