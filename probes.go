package probes

import (
	"context"
	"errors"
)

// Probes содержит контракты для реализации проб
// Kubernetes: Liveness, Readiness и Startup.
type Probes interface {
	Liveness
	Readiness
	Startup
}

// DefaultProbes содержит обработчики Liveness-, Readiness-
// и Startup-проб Kubernetes по-умолчанию.
//
//	Смотри DefaultLiveness, DefaultReadiness и DefaultStartup
var DefaultProbes = NewSuccessProbes()

// SuccessProbes содержит реализации проб Kubernetes,
// которые возвращают всегда положительный результат.
//
//	Смотри SuccessLiveness, SuccessReadiness, SuccessStartup
type SuccessProbes struct {
	successLiveness  Liveness
	successReadiness Readiness
	successStartup   Startup
}

func (probes *SuccessProbes) Liveness(ctx context.Context) (Result, error) {
	return probes.successLiveness.Liveness(ctx)
}

func (probes *SuccessProbes) Readiness(ctx context.Context) (Result, error) {
	return probes.successReadiness.Readiness(ctx)
}

func (probes *SuccessProbes) Startup(ctx context.Context) (Result, error) {
	return probes.successStartup.Startup(ctx)
}

// NewSuccessProbes инициализирует пробы Kubernetes,
// которые возвращают всегда положительный результат.
func NewSuccessProbes() *SuccessProbes {
	return &SuccessProbes{
		successLiveness:  DefaultLiveness,
		successReadiness: DefaultReadiness,
		successStartup:   DefaultStartup,
	}
}

// ErrUnsupportedResult указывает, что используемый
// результат в ответе эндпоинта не поддерживается.
//
//	Смотри Result
var ErrUnsupportedResult = errors.New("probes: unsupported result")

var results = map[Result]string{
	Success: "success",
	Warning: "warning",
	Failure: "failure",
}

// Result определяет результат, возвращаемый
// эндпоинтами Kubernetes: Liveness, Readiness
// и Startup.
//
//	Поддерживаемые статусы:
//	Success
//	Warning
//	Failure
type Result uint8

const (
	// Success указывает на успешный ответ эндпоинта.
	Success Result = iota

	// Warning логически означает успешный ответ, но
	// с дополнительной отладочной информацией.
	Warning

	// Failure указывает на ответ с ошибкой от эндпоинта.
	Failure
)

func (r Result) String() string {
	result, found := results[r]
	if found {
		return result
	}

	return ""
}

func (r Result) Validate() error {
	_, found := results[r]
	if found {
		return nil
	}

	return ErrUnsupportedResult
}

// IsSuccess возвращает true, если результат
// равен Success.
func (r Result) IsSuccess() bool {
	return r.is(Success)
}

// IsWarning возвращает true, если результат
// равен Warning.
func (r Result) IsWarning() bool {
	return r.is(Warning)
}

// IsFailure возвращает true, если результат
// равен Failure.
func (r Result) IsFailure() bool {
	return r.is(Failure)
}

func (r Result) is(other Result) bool {
	return r == other
}
