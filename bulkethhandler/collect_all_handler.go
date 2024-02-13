package bulkethhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Roman2K/txstreet-bulk-eth-api/collection"
	"github.com/Roman2K/txstreet-bulk-eth-api/limits"
)

const defaultWriteTimeout = 30 * time.Second

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
		err = fmt.Errorf("Failed to decode request body: %w", err)
		sendError(w, r, err)
		return
	}

	slog.Debug("Received collect request", "req", contractsReq)

	var cancel context.CancelFunc
	r, cancel = setRequestTimeout(r)
	defer cancel()

	resultsCollector := collection.LimitCollector[Input, Result]{
		Limiter: h.limiter,
		CollectFunc: func(elem Input) (Result, error) {
			return h.collectAdapter.collectFunc(r.Context(), elem)
		},
	}

	results, err := resultsCollector.CollectAll(h.collectAdapter.getInputs(contractsReq))
	if err != nil {
		sendError(w, r, fmt.Errorf("Failed to collect results: %w", err))
		return
	}

	slog.Debug("Sending results", "results", results)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		sendError(w, r, fmt.Errorf("Failed to JSON-encode results: %w", err))
		return
	}
}

var _ http.Handler = collectAllHandler[int, int, int]{}

func setRequestTimeout(r *http.Request) (*http.Request, context.CancelFunc) {
	var (
		logger                   = requestLogger(r)
		requestCtx               = r.Context()
		timeout    time.Duration = defaultWriteTimeout
	)

	if server, ok := requestCtx.Value(http.ServerContextKey).(*http.Server); ok {
		if t := server.WriteTimeout; t > 0 && t < timeout {
			timeout = t
		}
	}

	logger.Debug("Setting request timeout", "timeout", timeout)

	var cancel context.CancelFunc
	requestCtx, cancel = context.WithTimeout(requestCtx, timeout)

	return r.WithContext(requestCtx), cancel
}
