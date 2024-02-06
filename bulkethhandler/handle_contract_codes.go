package bulkethhandler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Roman2K/bulk-eth-api/collection"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type contractCodesRequest struct {
	Contracts []common.Address `json:"contracts"`
}

type contractCodeResult struct {
	Contract common.Address `json:"contract"`
	Code     string         `json:"code"`
}

func (h handler) handleContractCodes(w http.ResponseWriter, r *http.Request) {
	contractsReq, err := requestBodyDecoder[contractCodesRequest]{w, r}.Decode()
	if err != nil {
		sendError(
			w,
			fmt.Errorf("Failed to decode contract codes request body: %w", err),
		)
		return
	}

	slog.Debug("Received contract codes request", "req", contractsReq)

	resultsCollector := collection.LimitCollector[common.Address, contractCodeResult]{
		Limiter: h.limiter,
		CollectFunc: func(account common.Address) (result contractCodeResult, err error) {
			result.Contract = account

			var code []byte
			code, err = h.ethClient.CodeAt(h.ctx, account, nil)
			result.Code = hexutil.Encode(code)

			if err != nil {
				err = fmt.Errorf(
					"Failed to get contract code of account %s: %w", account, err,
				)
			}

			return
		},
	}

	results, err := resultsCollector.CollectAll(contractsReq.Contracts)
	if err != nil {
		sendError(w, fmt.Errorf("Failed to collect contract codes: %w", err))
		return
	}

	slog.Debug("Returning contract codes results", "results", results)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		sendError(w, fmt.Errorf("Failed to JSON-encode results: %w", err))
		return
	}
}
