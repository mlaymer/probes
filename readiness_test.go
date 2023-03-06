package probes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessReadiness_Readiness(t *testing.T) {
	t.Parallel()

	// Arrange.
	readiness := SuccessReadiness{}

	// Act.
	result, err := readiness.Readiness(context.Background())

	// Assert.
	assert.NoError(t, err)

	assert.Equal(t, Success, result)
}
