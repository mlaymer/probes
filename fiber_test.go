package probes

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

const fiberAddress = ":9001"

var dummyFiberError = errors.New("fiber: dummy error")

func TestFiber(t *testing.T) {
	t.Parallel()

	// Arrange.
	app := fiber.New()

	// Act.
	server := Fiber(app)

	// Assert.
	require.NotNil(t, server)

	assert.Nil(t, server.probes)
	assert.Equal(t, server.app, app)
}

func TestFiberServer_Probes(t *testing.T) {
	t.Parallel()

	// Arrange.
	app := fiber.New()

	server := &FiberServer{
		app:    app,
		probes: nil,
	}

	// Act.
	server = server.Probes(DefaultProbes)

	// Act.
	assert.Equal(t, DefaultProbes, server.probes)
}

func TestFiberServer_Start(t *testing.T) {
	// Arrange.
	app := fiber.New()

	server := Fiber(app).Probes(DefaultProbes)

	// Act.
	go func() {
		err := server.Start(fiberAddress)

		assert.NoError(t, err)
	}()

	time.Sleep(time.Second)

	{
		response, err := app.Test(httptest.NewRequest(fiber.MethodGet, DefaultLivenessPath, nil))

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	}

	{
		response, err := app.Test(httptest.NewRequest(fiber.MethodGet, DefaultReadinessPath, nil))

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	}

	{
		response, err := app.Test(httptest.NewRequest(fiber.MethodGet, DefaultStartupPath, nil))

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	}
}

type testSuccessLiveness struct{}

func (t testSuccessLiveness) Liveness(context.Context) (Result, error) {
	return Success, nil
}

type testWarningLivenessErrorless struct{}

func (t testWarningLivenessErrorless) Liveness(context.Context) (Result, error) {
	return Warning, nil
}

type testWarningLivenessWithError struct{}

func (t testWarningLivenessWithError) Liveness(context.Context) (Result, error) {
	return Warning, dummyFiberError
}

type testFailureLivenessErrorless struct{}

func (t testFailureLivenessErrorless) Liveness(context.Context) (Result, error) {
	return Failure, nil
}

type testFailureLivenessWithError struct{}

func (t testFailureLivenessWithError) Liveness(context.Context) (Result, error) {
	return Failure, dummyFiberError
}

func TestFiberLiveness_Liveness(t *testing.T) {
	t.Parallel()

	// Arrange.
	type responseChecker func(t *testing.T, ctx *fiber.Ctx)

	tests := []struct {
		name        string
		handler     *FiberLiveness
		ctx         *fiber.Ctx
		expectedErr error
		checker     responseChecker
	}{
		{
			name:    `Обработчик ответил "Success"`,
			handler: &FiberLiveness{probe: testSuccessLiveness{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Warning", но без ошибки`,
			handler: &FiberLiveness{probe: testWarningLivenessErrorless{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Warning" с ошибкой`,
			handler: &FiberLiveness{probe: testWarningLivenessWithError{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Equal(t, dummyFiberError.Error(), string(ctx.Response().Body()))
			},
		},
		{
			name:    `Обработчик ответил "Failure", но без ошибки`,
			handler: &FiberLiveness{probe: testFailureLivenessErrorless{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Failure" с ошибкой`,
			handler: &FiberLiveness{probe: testFailureLivenessWithError{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
				assert.Equal(t, dummyFiberError.Error(), string(ctx.Response().Body()))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actualErr := test.handler.Liveness(test.ctx)

			// Assert.
			assert.Equal(t, test.expectedErr, actualErr)

			test.checker(t, test.ctx)
		})
	}
}

func TestNewFiberLiveness(t *testing.T) {
	t.Parallel()

	// Act.
	actualLiveness := NewFiberLiveness(DefaultLiveness)

	// Assert.
	assert.Equal(t, DefaultLiveness, actualLiveness.probe)
}

type testSuccessReadiness struct{}

func (t testSuccessReadiness) Readiness(context.Context) (Result, error) {
	return Success, nil
}

type testWarningReadinessErrorless struct{}

func (t testWarningReadinessErrorless) Readiness(context.Context) (Result, error) {
	return Warning, nil
}

type testWarningReadinessWithError struct{}

func (t testWarningReadinessWithError) Readiness(context.Context) (Result, error) {
	return Warning, dummyFiberError
}

type testFailureReadinessErrorless struct{}

func (t testFailureReadinessErrorless) Readiness(context.Context) (Result, error) {
	return Failure, nil
}

type testFailureReadinessWithError struct{}

func (t testFailureReadinessWithError) Readiness(context.Context) (Result, error) {
	return Failure, dummyFiberError
}

func TestFiberReadiness_Readiness(t *testing.T) {
	t.Parallel()

	// Arrange.
	type responseChecker func(t *testing.T, ctx *fiber.Ctx)

	tests := []struct {
		name        string
		handler     *FiberReadiness
		ctx         *fiber.Ctx
		expectedErr error
		checker     responseChecker
	}{
		{
			name:    `Обработчик ответил "Success"`,
			handler: &FiberReadiness{probe: testSuccessReadiness{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Warning", но без ошибки`,
			handler: &FiberReadiness{probe: testWarningReadinessErrorless{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Warning" с ошибкой`,
			handler: &FiberReadiness{probe: testWarningReadinessWithError{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Equal(t, dummyFiberError.Error(), string(ctx.Response().Body()))
			},
		},
		{
			name:    `Обработчик ответил "Failure", но без ошибки`,
			handler: &FiberReadiness{probe: testFailureReadinessErrorless{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Failure" с ошибкой`,
			handler: &FiberReadiness{probe: testFailureReadinessWithError{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
				assert.Equal(t, dummyFiberError.Error(), string(ctx.Response().Body()))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actualErr := test.handler.Readiness(test.ctx)

			// Assert.
			assert.Equal(t, test.expectedErr, actualErr)

			test.checker(t, test.ctx)
		})
	}
}

func TestNewFiberReadiness(t *testing.T) {
	t.Parallel()

	// Act.
	actualReadiness := NewFiberReadiness(DefaultReadiness)

	// Assert.
	assert.Equal(t, DefaultReadiness, actualReadiness.probe)
}

type testSuccessStartup struct{}

func (t testSuccessStartup) Startup(context.Context) (Result, error) {
	return Success, nil
}

type testWarningStartupErrorless struct{}

func (t testWarningStartupErrorless) Startup(context.Context) (Result, error) {
	return Warning, nil
}

type testWarningStartupWithError struct{}

func (t testWarningStartupWithError) Startup(context.Context) (Result, error) {
	return Warning, dummyFiberError
}

type testFailureStartupErrorless struct{}

func (t testFailureStartupErrorless) Startup(context.Context) (Result, error) {
	return Failure, nil
}

type testFailureStartupWithError struct{}

func (t testFailureStartupWithError) Startup(context.Context) (Result, error) {
	return Failure, dummyFiberError
}

func TestFiberStartup_Startup(t *testing.T) {
	t.Parallel()

	// Arrange.
	type responseChecker func(t *testing.T, ctx *fiber.Ctx)

	tests := []struct {
		name        string
		handler     *FiberStartup
		ctx         *fiber.Ctx
		expectedErr error
		checker     responseChecker
	}{
		{
			name:    `Обработчик ответил "Success"`,
			handler: &FiberStartup{probe: testSuccessStartup{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Warning", но без ошибки`,
			handler: &FiberStartup{probe: testWarningStartupErrorless{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Warning" с ошибкой`,
			handler: &FiberStartup{probe: testWarningStartupWithError{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
				assert.Equal(t, dummyFiberError.Error(), string(ctx.Response().Body()))
			},
		},
		{
			name:    `Обработчик ответил "Failure", но без ошибки`,
			handler: &FiberStartup{probe: testFailureStartupErrorless{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
				assert.Nil(t, ctx.Response().Body())
			},
		},
		{
			name:    `Обработчик ответил "Failure" с ошибкой`,
			handler: &FiberStartup{probe: testFailureStartupWithError{}},
			ctx: func() *fiber.Ctx {
				app := fiber.New()

				ctx := app.AcquireCtx(new(fasthttp.RequestCtx))

				return ctx
			}(),
			expectedErr: nil,
			checker: func(t *testing.T, ctx *fiber.Ctx) {
				assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
				assert.Equal(t, dummyFiberError.Error(), string(ctx.Response().Body()))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act.
			actualErr := test.handler.Startup(test.ctx)

			// Assert.
			assert.Equal(t, test.expectedErr, actualErr)

			test.checker(t, test.ctx)
		})
	}
}

func TestNewFiberStartup(t *testing.T) {
	t.Parallel()

	// Act.
	actualStartup := NewFiberStartup(DefaultStartup)

	// Assert.
	assert.Equal(t, DefaultStartup, actualStartup.probe)
}
