package probes

import (
	"github.com/gofiber/fiber/v2"
)

// Fiber подготавливает REST-сервер для приёма
// Liveness-, Readiness- и Startup-проб.
//
// После вызова метода следует вызвать метод
// FiberServer.Probes, который выполнит инициализацию
// REST-эндпоинтов.
func Fiber(app *fiber.App) *FiberServer {
	return &FiberServer{app: app}
}

// FiberServer содержит REST-эндпоинты для проб
// Kubernetes: Liveness, Readiness и Startup.
//
// Для инициализации сервера необходим сначала
// вызвать метод Fiber с необходимыми параметрами.
type FiberServer struct {
	app *fiber.App

	probes Probes
}

// Probes инициализирует REST-эндпоинты для Liveness-,
// Readiness- и Startup-проб Kubernetes.
//
//	Liveness-эндпоинт доступен по пути /liveness
//	Readiness-эндпоинт доступен по пути /readiness
//	Startup-эндпоинт доступен по пути /startup
func (server *FiberServer) Probes(probes Probes) *FiberServer {
	server.app.Get(DefaultLivenessPath, NewFiberLiveness(probes).Liveness)
	server.app.Get(DefaultReadinessPath, NewFiberReadiness(probes).Readiness)
	server.app.Get(DefaultStartupPath, NewFiberStartup(probes).Startup)

	server.probes = probes

	return server
}

// Start запускает REST-сервер с пробами Kubernetes
// по указанному адресу.
func (server *FiberServer) Start(address string) error {
	return server.app.Listen(address)
}

// DefaultFiberLiveness содержит инициализированный
// HTTP-обработчик Liveness-запроса Kubernetes.
//
// Обработчик всегда возвращает HTTP 200 OK.
var DefaultFiberLiveness = NewFiberLiveness(DefaultLiveness)

// FiberLiveness обрабатывает HTTP-запросы Liveness
// Kubernetes.
//
// Обработчику необходима реализация Liveness.
//
// Для инициализации необходимо использовать метод
// NewFiberLiveness.
// Реализация по умолчанию – DefaultFiberLiveness.
type FiberLiveness struct {
	probe Liveness
}

// Liveness обрабатывает HTTP-запрос Liveness от
// Kubernetes.
//
// Если Liveness возвращает Success или Warning,
// то обработчик возвращает HTTP 200 OK с опциональным
// телом, если Failure – HTTP 500 Internal Server Error
// с опциональным телом.
func (handler FiberLiveness) Liveness(ctx *fiber.Ctx) error {
	result, err := handler.probe.Liveness(ctx.Context())

	if result.IsSuccess() || result.IsWarning() {
		ctx.Status(fiber.StatusOK)
	}

	if result.IsFailure() {
		ctx.Status(fiber.StatusInternalServerError)
	}

	return sendFiberError(ctx, err)
}

// NewFiberLiveness инициализирует HTTP-обработчик
// Liveness-запросов Kubernetes на Fiber.
func NewFiberLiveness(probe Liveness) FiberLiveness {
	return FiberLiveness{probe: probe}
}

// DefaultFiberReadiness содержит инициализированный
// HTTP-обработчик Readiness-запроса Kubernetes.
//
// Обработчик всегда возвращает HTTP 200 OK.
var DefaultFiberReadiness = NewFiberReadiness(DefaultReadiness)

// FiberReadiness обрабатывает HTTP-запросы Readiness
// Kubernetes.
//
// Обработчику необходима реализация Readiness.
//
// Для инициализации необходимо использовать метод
// NewFiberReadiness.
// Реализация по умолчанию – DefaultFiberReadiness.
type FiberReadiness struct {
	probe Readiness
}

// Readiness обрабатывает HTTP-запрос Readiness от
// Kubernetes.
//
// Если Readiness возвращает Success или Warning,
// то обработчик возвращает HTTP 200 OK с опциональным
// телом, если Failure – HTTP 500 Internal Server Error
// с опциональным телом.
func (handler FiberReadiness) Readiness(ctx *fiber.Ctx) error {
	result, err := handler.probe.Readiness(ctx.Context())

	if result.IsSuccess() || result.IsWarning() {
		ctx.Status(fiber.StatusOK)
	}

	if result.IsFailure() {
		ctx.Status(fiber.StatusInternalServerError)
	}

	return sendFiberError(ctx, err)
}

// NewFiberReadiness инициализирует HTTP-обработчик
// Readiness-запросов Kubernetes на Fiber.
func NewFiberReadiness(probe Readiness) FiberReadiness {
	return FiberReadiness{probe: probe}
}

// DefaultFiberStartup содержит инициализированный
// HTTP-обработчик Startup-запроса Kubernetes.
//
// Обработчик всегда возвращает HTTP 200 OK.
var DefaultFiberStartup = NewFiberStartup(DefaultStartup)

// FiberStartup обрабатывает HTTP-запросы Startup
// Kubernetes.
//
// Обработчику необходима реализация Startup.
//
// Для инициализации необходимо использовать метод
// NewFiberStartup.
// Реализация по умолчанию – DefaultFiberStartup.
type FiberStartup struct {
	probe Startup
}

// Startup обрабатывает HTTP-запрос Startup от
// Kubernetes.
//
// Если Startup возвращает Success или Warning,
// то обработчик возвращает HTTP 200 OK с опциональным
// телом, если Failure – HTTP 500 Internal Server Error
// с опциональным телом.
func (handler FiberStartup) Startup(ctx *fiber.Ctx) error {
	result, err := handler.probe.Startup(ctx.Context())

	if result.IsSuccess() || result.IsWarning() {
		ctx.Status(fiber.StatusOK)
	}

	if result.IsFailure() {
		ctx.Status(fiber.StatusInternalServerError)
	}

	return sendFiberError(ctx, err)
}

// NewFiberStartup инициализирует HTTP-обработчик
// Startup-запросов Kubernetes на Fiber.
func NewFiberStartup(probe Startup) FiberStartup {
	return FiberStartup{probe: probe}
}

func sendFiberError(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return ctx.Send(nil)
	}

	return ctx.SendString(err.Error())
}
