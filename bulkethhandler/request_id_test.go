package bulkethhandler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithRequestId(t *testing.T) {
	ctx := context.Background()

	requestId, ok := contextRequestId(ctx)

	assert.False(t, ok)
	assert.Nil(t, requestId)

	requestId = newRequestId()
	requestIdCtx := withRequestId(ctx, requestId)

	gotRequestId, ok := contextRequestId(requestIdCtx)

	assert.True(t, ok)
	assert.Equal(t, requestId, gotRequestId)
}
