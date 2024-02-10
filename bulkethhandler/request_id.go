package bulkethhandler

import (
	"context"
	"encoding"
	"fmt"

	"github.com/segmentio/ksuid"
)

type requestId interface {
	fmt.Stringer
	encoding.TextMarshaler
}

func newRequestId() requestId {
	return ksuid.New()
}

type requestIdKey struct{}

func withRequestId(ctx context.Context, reqId requestId) context.Context {
	return context.WithValue(ctx, requestIdKey{}, reqId)
}

func contextRequestId(ctx context.Context) (requestId, bool) {
	reqId, ok := ctx.Value(requestIdKey{}).(requestId)
	return reqId, ok
}
