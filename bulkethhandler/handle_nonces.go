package bulkethhandler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/collection"
	"github.com/ethereum/go-ethereum/common"
)

type noncesRequest struct {
	Accounts []common.Address `json:"accounts"`
}

type nonceResult struct {
	Account common.Address `json:"account"`
	Count   uint64         `json:"count"`
}

func (h handler) handleNonces(w http.ResponseWriter, r *http.Request) {
	noncesReq, err := requestBodyDecoder[noncesRequest]{w, r}.Decode()
	if err != nil {
		sendError(w, fmt.Errorf("Failed to decode nonces request body: %w", err))
		return
	}

	slog.Debug("Received nonces request", "req", noncesReq)

	resultsCollector := collection.LimitCollector[common.Address, nonceResult]{
		Limiter: h.limiter,
		CollectFunc: func(account common.Address) (result nonceResult, err error) {
			result.Account = account
			result.Count, err = h.ethClient.NonceAt(h.ctx, account, nil)

			if err != nil {
				err = fmt.Errorf("Failed to get nonce of account %s: %w", account, err)
			}

			return
		},
	}

	results, err := resultsCollector.CollectAll(noncesReq.Accounts)
	if err != nil {
		sendError(w, fmt.Errorf("Failed to collect nonce results: %w", err))
		return
	}

	slog.Debug("Returning nonces results", "results", results)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		sendError(w, fmt.Errorf("Failed to JSON-encode results: %w", err))
		return
	}
}
