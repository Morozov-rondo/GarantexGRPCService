### Сервис получения курса USDT

## Описание

Данный сервис предоставляет GRPC-методы для получения текущего курса USDT с биржи Garantex. Сервис хранит полученные данные в базе данных PostgreSQL и предоставляет healthcheck-метод для проверки работоспособности.

## Функциональные возможности

1. GetRates: GRPC-метод для получения текущего курса USDT (ask, bid и метка времени).
2. Healthcheck: GRPC-метод для проверки работоспособности сервиса.
3. Graceful shutdown: Корректное завершение работы сервиса.

## Технологии и зависимости

Язык программирования: Go
База данных: PostgreSQL
Сборка и запуск: Docker, Docker Compose
Статический анализ кода: golangci-lint

## Запуск
Для запуска сервиса используйте Docker Compose:

```
docker compose up
```

Сервис будет доступен по адресу localhost:8080.

## Тестирование
Для запуска тестов используйте команду:

```
make test
```
Команда make test запускает все unit-тесты в проекте с флагами -v -cover, которые показывают подробную информацию о тестах и процент покрытия кода.

## Сборка и запуск в Docker
Для сборки Docker-образа используйте команду:

```
make docker-build
```

Команда make docker-build собирает Docker-образ с тегом exchange:dev.

Для запуска сервиса в Docker используйте команду:

```
make run
```

Команда make run запускает сервис в Docker Compose.

## Статический анализ кода

Для запуска статического анализа кода используйте команду:

```
make lint
```

Команда make lint запускает golangci-lint, который проверяет весь код в проекте на наличие ошибок и стилевых нарушений.
