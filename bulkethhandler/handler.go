package bulkethhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/collection"
	"github.com/Roman2K/bulk-eth-api/limits"
	"github.com/ethereum/go-ethereum/common"
)

type Handler struct {
	Context   context.Context
	EthClient EthClient
	Limiter   limits.Limiter
}

type noncesRequest struct {
	Accounts []common.Address `json:"accounts"`
}

type nonceResult struct {
	Account common.Address `json:"account"`
	Count   uint64         `json:"count"`
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	noncesReq, err := requestBodyDecoder[noncesRequest]{w, r}.Decode()
	if err != nil {
		sendError(w, fmt.Errorf("Failed to decode nonces request body: %w", err))
		return
	}

	log.Printf("Received nonces request: %v\n", noncesReq)

	resultsCollector := collection.LimitCollector[common.Address, nonceResult]{
		Limiter: h.Limiter,
		CollectFunc: func(account common.Address) (result nonceResult, err error) {
			result.Account = account
			result.Count, err = h.EthClient.NonceAt(h.Context, account, nil)

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

	log.Printf("Returning nonces results: %v\n", results)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		sendError(w, fmt.Errorf("Failed to JSON-encode results: %w", err))
		return
	}
}

func sendError(w http.ResponseWriter, err error) {
	log.Printf("Error while serving request: %s", err)
	http.Error(w, "Internal error", 500)
}
