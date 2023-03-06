package probes

import "context"

// Liveness оборачивает метод Liveness, который
// обрабатывает Liveness-запрос Kubernetes, используемый
// для определения нужно ли перезапускать контейнер.
//
// Во время работы приложения накапливаются ошибки, которые
// исправляются только после перезапуска контейнера.
//
// Метод возвращает Result и ошибку, которая содержит
// отладочную информацию.
type Liveness interface {
	Liveness(context.Context) (Result, error)
}

// DefaultLiveness содержит обработчик Liveness-проб
// Kubernetes по умолчанию.
//
// Обработчик по умолчанию – SuccessLiveness.
var DefaultLiveness = SuccessLiveness{}

// SuccessLiveness реализует обработку Liveness-пробы
// Kubernetes с результатом Success.
//
// Обработчик при любых обстоятельствах возвращает
// именно этот результат, ошибка в ответе отсутствует.
type SuccessLiveness struct{}

func (probe SuccessLiveness) Liveness(context.Context) (Result, error) {
	return Success, nil
}
