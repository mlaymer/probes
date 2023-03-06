package probes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult_String(t *testing.T) {
	t.Parallel()

	// Arrange.
	tests := []struct {
		name           string
		result         Result
		expectedString string
	}{
		{
			name:           `Результат "Success"`,
			result:         Success,
			expectedString: "success",
		},
		{
			name:           `Результат "Warning"`,
			result:         Warning,
			expectedString: "warning",
		},
		{
			name:           `Результат "Failure"`,
			result:         Failure,
			expectedString: "failure",
		},
		{
			name:           "Неподдерживаемый результат",
			result:         100,
			expectedString: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actualString := test.result.String()

			// Assert.
			assert.Equal(t, test.expectedString, actualString)
		})
	}
}

func TestResult_Validate(t *testing.T) {
	t.Parallel()

	// Arrange.
	tests := []struct {
		name        string
		result      Result
		expectedErr error
	}{
		{
			name:        `Результат "Success"`,
			result:      Success,
			expectedErr: nil,
		},
		{
			name:        `Результат "Warning"`,
			result:      Warning,
			expectedErr: nil,
		},
		{
			name:        `Результат "Failure"`,
			result:      Failure,
			expectedErr: nil,
		},
		{
			name:        "Неподдерживаемый результат",
			result:      100,
			expectedErr: ErrUnsupportedResult,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actualErr := test.result.Validate()

			// Assert.
			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

func TestResult_IsSuccess(t *testing.T) {
	t.Parallel()

	// Arrange.
	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{
			name:     `Результат "Success"`,
			result:   Success,
			expected: true,
		},
		{
			name:     `Результат "Warning"`,
			result:   Warning,
			expected: false,
		},
		{
			name:     `Результат "Failure"`,
			result:   Failure,
			expected: false,
		},
		{
			name:     "Неподдерживаемый результат",
			result:   100,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actual := test.result.IsSuccess()

			// Assert.
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestResult_IsWarning(t *testing.T) {
	t.Parallel()

	// Arrange.
	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{
			name:     `Результат "Success"`,
			result:   Success,
			expected: false,
		},
		{
			name:     `Результат "Warning"`,
			result:   Warning,
			expected: true,
		},
		{
			name:     `Результат "Failure"`,
			result:   Failure,
			expected: false,
		},
		{
			name:     "Неподдерживаемый результат",
			result:   100,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actual := test.result.IsWarning()

			// Assert.
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestResult_IsFailure(t *testing.T) {
	t.Parallel()

	// Arrange.
	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{
			name:     `Результат "Success"`,
			result:   Success,
			expected: false,
		},
		{
			name:     `Результат "Warning"`,
			result:   Warning,
			expected: false,
		},
		{
			name:     `Результат "Failure"`,
			result:   Failure,
			expected: true,
		},
		{
			name:     "Неподдерживаемый результат",
			result:   100,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actual := test.result.IsFailure()

			// Assert.
			assert.Equal(t, test.expected, actual)
		})
	}
}
