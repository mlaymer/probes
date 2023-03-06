# Golang Kubernetes Probes Kit

Комплект разработчика для обслуживания проб Kubernetes.

## Probes Kit: Введение для разработчиков

> Контракты для Liveness-, Readiness- и Startup-проб Kubernetes.

> Интеграция с [Fiber](https://github.com/gofiber/fiber).

## Использование

### Поддержка проб Kubernetes

Kubernetes Kit предоставляет контракт для реализации обработчиков запросов проб Kubernetes: Liveness, Readiness и
Startup.

Для создания собственных обработчиков необходимо реализовать `probes.Probes`, который включает `probes.Liveness`,
`probes.Readiness` и `probes.Startup`.

Пример контракта `probes.Liveness`:

```go
package probes

type Liveness interface {
	Liveness(context.Context) (Result, error)
}
```

Все контракты обработчиков запросов проб оперируют `probes.Result`, предоставляющий несколько значений результата
обработки запроса. Запрос может быть выполнен успешно, успешно с отладочной информацией и с ошибкой.

### Интеграция с Fiber

Probes Kit интегрирован с [Fiber](https://github.com/gofiber/fiber).

Для всех проб существуют HTTP-обработчики Fiber: `probes.FiberLiveness`, `probes.FiberReadiness`, `probes.FiberStartup`.

#### Инициализация и запуск REST-сервера

```go
package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mlaymer/probes"
)

func main() {
	app := fiber.New()

	server := probes.Fiber(app).Probes(probes.DefaultProbes)

	log.Println(server.Start(probes.DefaultServerAddress))
}
```
