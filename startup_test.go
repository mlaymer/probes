package probes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessStartup_Startup(t *testing.T) {
	t.Parallel()

	// Arrange.
	startup := SuccessStartup{}

	// Act.
	result, err := startup.Startup(context.Background())

	// Assert.
	assert.NoError(t, err)

	assert.Equal(t, Success, result)
}
