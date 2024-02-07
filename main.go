package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Roman2K/bulk-eth-api/bulkethhandler"
	"github.com/Roman2K/bulk-eth-api/limits"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("vim-go")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var opts options
	opts.flagSet().Parse(os.Args[1:])
	slog.Info(
		"Command-line args parsed",
		"listenAddr", opts.listenAddr,
		"ethUrl", opts.ethUrl,
		"ethConcurrency", opts.ethConcurrency,
		"logLevel", opts.logLevel,
	)

	setLogger(opts.logLevel)
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, opts.ethUrl)
	if err != nil {
		return fmt.Errorf("Failed to instantiate Ethereum client: %w", err)
	}

	limiter := limits.NewLimiter(opts.ethConcurrency)
	handler := bulkethhandler.NewHandler(ctx, client, limiter)

	slog.Info("Listening", "addr", opts.listenAddr)
	return http.ListenAndServe(opts.listenAddr, handler)
}

func setLogger(level slog.Level) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stderr,
				&slog.HandlerOptions{Level: level},
			),
		),
	)
}
