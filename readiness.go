package probes

import "context"

// Readiness оборачивает метод Readiness, который
// обрабатывает Readiness-запрос Kubernetes, используемый
// для определения готов ли контейнер принимать трафик.
//
// Сигнал Readiness используется для включения пода в
// балансировщик нагрузки.
//
// Метод возвращает Result и ошибку, которая содержит
// отладочную информацию.
type Readiness interface {
	Readiness(context.Context) (Result, error)
}

// DefaultReadiness содержит обработчик Readiness-проб
// Kubernetes по умолчанию.
//
// Обработчик по умолчанию – SuccessReadiness.
var DefaultReadiness = SuccessReadiness{}

// SuccessReadiness реализует обработку Readiness-пробы
// Kubernetes с результатом Success.
//
// Обработчик при любых обстоятельствах возвращает
// именно этот результат, ошибка в ответе отсутствует.
type SuccessReadiness struct{}

func (s SuccessReadiness) Readiness(context.Context) (Result, error) {
	return Success, nil
}
