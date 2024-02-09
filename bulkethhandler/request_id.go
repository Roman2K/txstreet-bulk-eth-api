package bulkethhandler

import (
	"context"
	"encoding"
	"errors"
	"fmt"
)

var errMissingRequestId = errors.New("Missing request ID in context")

type requestId interface {
	fmt.Stringer
	encoding.TextMarshaler
}

type requestKey string

var requestIdKey = requestKey("requestId")

func setRequestId(ctx context.Context, reqId requestId) context.Context {
	return context.WithValue(ctx, requestIdKey, reqId)
}

func getRequestId(ctx context.Context) (requestId, bool) {
	reqId, ok := ctx.Value(requestIdKey).(requestId)
	return reqId, ok
}
