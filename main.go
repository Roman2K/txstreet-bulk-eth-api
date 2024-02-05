package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Roman2K/bulk-eth-api/bulkethhandler"
	"github.com/Roman2K/bulk-eth-api/limits"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("vim-go")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const listenAddr = ":8081"

func run() error {
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, "/mnt/lfs/geth/data/ipc")
	if err != nil {
		return fmt.Errorf("Failed to instantiate Ethereum client: %w", err)
	}

	limiter := limits.NewLimiter(2)

	handler := bulkethhandler.Handler{
		Context:   ctx,
		EthClient: client,
		Limiter:   limiter,
	}

	log.Printf("Listening on %s\n", listenAddr)

	return http.ListenAndServe(listenAddr, handler)

	// /nonces

	account := common.HexToAddress("0x7d3fc733f2a39af39cb0af950598950c82925749")
	nonce, err := client.NonceAt(ctx, account, nil)
	if err != nil {
		return fmt.Errorf("Failed to get nonce of account %s: %w", account, err)
	}

	fmt.Printf("nonce: %d\n", nonce)

	// /contract-codes

	// USDC
	account = common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	code, err := client.CodeAt(ctx, account, nil)
	if err != nil {
		return fmt.Errorf("Failed to get code of account %s: %w", account, err)
	}

	fmt.Printf("code:\n%s\n", hexutil.Encode(code))

	// /transaction-receipts

	txHash := common.HexToHash("0x6f9f225f36d1a6b50d39e7d9fe3a2327adf29833bfbacbe76d7a292c4fa087f7")
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return fmt.Errorf("Failed to get receipt of transaction %s: %w", txHash, err)
	}

	logs := make([]string, 0, len(receipt.Logs))
	for _, log := range receipt.Logs {
		data, err := log.MarshalJSON()
		if err != nil {
			return err
		}

		logs = append(logs, string(data))
	}
	fmt.Printf("logs: %v\n", strings.Join(logs, "\n"))

	return nil
}
