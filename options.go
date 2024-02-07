package main

import (
	"flag"
	"log/slog"
	"strings"
)

type options struct {
	listenAddr     string
	ethUrl         string
	ethConcurrency int
	logLevel       slog.Level
}

func (opts *options) flagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("bulk-eth-api", flag.ExitOnError)

	flagSet.StringVar(
		&opts.listenAddr, "addr", ":8081",
		"Listen interface and port",
	)
	flagSet.StringVar(
		&opts.ethUrl, "eth", "http://localhost:8545",
		"RPC URL of Ethereum execution client",
	)
	flagSet.IntVar(
		&opts.ethConcurrency, "concurrency", 8,
		"Max concurrency of RPC calls",
	)

	opts.logLevel = slog.LevelDebug
	flagSet.Var(
		logFlagValue{&opts.logLevel}, "log",
		"Log level ("+strings.Join(logLevelNames(), ", ")+")",
	)

	return flagSet
}

type logFlagValue struct {
	level *slog.Level
}

func (val logFlagValue) String() string {
	if val.level == nil {
		return ""
	}

	return val.level.String()
}

func (val logFlagValue) Set(value string) error {
	return val.level.UnmarshalText([]byte(value))
}

var _ flag.Value = logFlagValue{}

func logLevelNames() []string {
	levels := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	}

	names := make([]string, len(levels))
	for i, level := range levels {
		names[i] = level.String()
	}

	return names
}
