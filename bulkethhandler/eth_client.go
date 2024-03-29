package bulkethhandler

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
)

// Subset of ethereum/go-ethereum.Client
// https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient
type EthClient interface {
	Client() *ethRpc.Client
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}
