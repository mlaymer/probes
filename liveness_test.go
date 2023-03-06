package probes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessLiveness_Liveness(t *testing.T) {
	t.Parallel()

	// Arrange.
	liveness := SuccessLiveness{}

	// Act.
	result, err := liveness.Liveness(context.Background())

	// Assert.
	assert.NoError(t, err)

	assert.Equal(t, Success, result)
}
