package probes

import "context"

// Startup оборачивает метод Startup, который
// обрабатывает Startup-метод Kubernetes, используемый
// для определения запущено ли приложение.
//
// Если Kubernetes настроен на использование этой пробы,
// то пробы Liveness и Readiness отключаются до получения
// результата от Startup.
// Это сделано во избежание уничтожения их Kubernetes до
// момента окончательного старта. Используется в приложениях
// с медленным запуском.
//
// Метод возвращает Result и ошибку, содержащую отладочную
// информацию.
type Startup interface {
	Startup(context.Context) (Result, error)
}

// DefaultStartup содержит обработчик Startup-проб
// Kubernetes по умолчанию.
//
// Обработчик по умолчанию – SuccessStartup.
var DefaultStartup = SuccessStartup{}

// SuccessStartup реализует обработку Startup-пробы
// Kubernetes с результатом Success.
//
// Обработчик при любых обстоятельствах возвращает
// именно этот результат, ошибка в ответе отсутствует.
type SuccessStartup struct{}

func (probe SuccessStartup) Startup(context.Context) (Result, error) {
	return Success, nil
}
