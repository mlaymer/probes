package probes

const (
	// DefaultServerAddress содержит путь по умолчанию
	// REST-сервера для обработчиков проб Kubernetes.
	DefaultServerAddress = "localhost:9000"

	// DefaultLivenessPath содержит путь по умолчанию
	// для HTTP-обработчика Liveness-запросов Kubernetes.
	DefaultLivenessPath = "/liveness"

	// DefaultReadinessPath содержит путь по умолчанию
	// для HTTP-обработчика Readiness-запросов Kubernetes.
	DefaultReadinessPath = "/readiness"

	// DefaultStartupPath содержит путь по умолчанию
	// для HTTP-обработчика Startup-запросов Kubernetes.
	DefaultStartupPath = "/startup"
)
