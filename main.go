package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	var cancel context.CancelCauseFunc
	ctx, cancel = stopOnSignal(ctx)

	client, err := ethclient.DialContext(ctx, opts.ethUrl)
	if err != nil {
		return fmt.Errorf("Failed to instantiate Ethereum client: %w", err)
	}

	limiter := limits.NewLimiter(opts.ethConcurrency)
	handler := bulkethhandler.NewHandler(ctx, client, limiter)

	server := &http.Server{
		Addr:           opts.listenAddr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	slog.Info("HTTP server listening", "addr", server.Addr)

	go func() {
		defer cancel(nil)

		<-ctx.Done()
		slog.Debug("Signal context cancelled", "error", ctx.Err(), "cause", context.Cause(ctx))

		ctx, timeoutCancel := context.WithTimeout(ctx, 10*time.Second)
		defer timeoutCancel()

		slog.Info("Shutting down HTTP server", "timeout", formatContextTimeout(ctx))

		server.Shutdown(ctx)
	}()

	if err = server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	slog.Info("HTTP server shut down")

	slog.Info("Waiting for pending tasks")
	limiter.Wait()

	return nil
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

func formatContextTimeout(ctx context.Context) string {
	if deadline, ok := ctx.Deadline(); ok {
		return deadline.Sub(time.Now()).String()
	}

	return "None"
}

func stopOnSignal(ctx context.Context) (context.Context, context.CancelCauseFunc) {
	sigContext, cancel := context.WithCancelCause(ctx)

	sigChan := make(chan os.Signal)
	go func() {
		defer cancel(nil)

		signal := <-sigChan
		slog.Info("Received signal", "signal", signal)
		cancel(fmt.Errorf("Received %s", signal))

		signal = <-sigChan
		slog.Warn("Received second signal, stopping", "signal", signal)
		os.Exit(2)
	}()

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	return sigContext, cancel
}
