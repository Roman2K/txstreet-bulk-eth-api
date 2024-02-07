package bulkethhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/collection"
	"github.com/Roman2K/bulk-eth-api/limits"
)

type collectAllHandler[Request, Input, Result any] struct {
	limiter        limits.Limiter
	collectAdapter collectAdapter[Request, Input, Result]
}

type collectAdapter[Request, Input, Result any] interface {
	getInputs(Request) []Input
	collectFunc(context.Context, Input) (Result, error)
}

func (h collectAllHandler[Request, Input, Result]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contractsReq, err := requestBodyDecoder[Request]{w, r}.Decode()
	if err != nil {
		sendError(
			w,
			fmt.Errorf("Failed to decode request body: %w", err),
		)
		return
	}

	slog.Debug("Received collect request", "req", contractsReq)

	resultsCollector := collection.LimitCollector[Input, Result]{
		Limiter: h.limiter,
		CollectFunc: func(elem Input) (Result, error) {
			return h.collectAdapter.collectFunc(r.Context(), elem)
		},
	}

	results, err := resultsCollector.CollectAll(h.collectAdapter.getInputs(contractsReq))
	if err != nil {
		sendError(w, fmt.Errorf("Failed to collect results: %w", err))
		return
	}

	slog.Debug("Sending results", "results", results)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		sendError(w, fmt.Errorf("Failed to JSON-encode results: %w", err))
		return
	}
}

var _ http.Handler = collectAllHandler[int, int, int]{}
