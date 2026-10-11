# Golang Todo App

REST API на Go для управления пользователями, задачами и статистикой по ним, плюс веб-интерфейс на одной HTML-странице.

## Стек

| Компонент        | Технология                                   |
|------------------|----------------------------------------------|
| Язык             | Go 1.24+                                     |
| HTTP             | стандартный `net/http`, без фреймворков      |
| База данных      | PostgreSQL 18, драйвер `jackc/pgx/v5`        |
| Миграции         | `golang-migrate` (запускается в Docker)      |
| Логи             | `go.uber.org/zap`                            |
| Конфигурация     | переменные окружения, `kelseyhightower/envconfig` |
| Валидация        | `go-playground/validator/v10`                |
| Документация API | Swagger (`swaggo/swag`)                      |
| Окружение        | Docker Compose                               |

---

## Быстрый старт

### Что нужно установить

- [Go](https://go.dev/dl/) 1.24 или новее
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (включает Docker Compose)
- `make`

**Windows:** `make` ставится командой `winget install ezwinports.make`. Все команды `make` запускайте из **Git Bash**, а не из PowerShell: Makefile использует bash. Docker Desktop должен быть запущен.

### Настройка

```bash
cp .env.example .env
```

Заполните в `.env` пустые значения `POSTGRES_USER`, `POSTGRES_PASSWORD` и `POSTGRES_DB`. Остальные переменные можно оставить как есть.

### Запуск для разработки (приложение на хосте)

```bash
make env-up            # поднять PostgreSQL
make migrate-up        # применить миграции
make env-port-forward  # открыть порт 5432 для приложения на хосте
make todoapp-run       # запустить приложение через go run
```

### Запуск целиком в Docker

```bash
make env-up
make migrate-up
make todoapp-deploy
```

`public/index.html` копируется в образ при сборке. После изменения страницы выполните `make todoapp-deploy` ещё раз.

### Адреса

- Веб-интерфейс: http://localhost:5050/
- Swagger UI: http://localhost:5050/swagger/
- API: http://localhost:5050/api/v1/

### Остановка

```bash
make todoapp-undeploy  # приложение в Docker
make env-port-close    # проброс порта БД
make env-down          # PostgreSQL
```

Данные БД хранятся в `out/pgdata`, логи — в `out/logs`. Обе папки переживают перезапуск. Удалить их можно командами `make env-cleanup` и `make logs-cleanup`.

---

## Команды Makefile

| Команда                     | Что делает                                              |
|-----------------------------|---------------------------------------------------------|
| `make help`                 | Список всех команд                                      |
| `make env-up` / `env-down`  | Запустить / остановить PostgreSQL                       |
| `make env-port-forward` / `env-port-close` | Открыть / закрыть порт `127.0.0.1:5432`  |
| `make env-cleanup`          | Удалить данные БД (`out/pgdata`)                        |
| `make logs-cleanup`         | Удалить логи (`out/logs`)                               |
| `make migrate-up` / `migrate-down` | Применить / откатить миграции                    |
| `make migrate-create seq=name` | Создать новую пару файлов миграции в `migrations/`   |
| `make seed` / `seed-clean`  | Добавить / удалить тестовых пользователей и задачи (`scripts/seed.sql`) |
| `make swagger-gen`          | Перегенерировать `docs/` по аннотациям в коде           |
| `make todoapp-run`          | Запустить приложение локально (`go run`)                |
| `make todoapp-deploy` / `todoapp-undeploy` | Собрать и запустить / остановить контейнер приложения |
| `make ps`                   | Показать запущенные сервисы                             |

---

## Переменные окружения

Задаются в `.env` (шаблон — `.env.example`).

| Переменная              | Обязательна | Описание                                    | Пример                           |
|-------------------------|:-----------:|---------------------------------------------|----------------------------------|
| `POSTGRES_USER`         | да          | Пользователь БД                             | `todoapp`                        |
| `POSTGRES_PASSWORD`     | да          | Пароль БД                                   | `secret`                         |
| `POSTGRES_DB`           | да          | Имя базы данных                             | `todoapp`                        |
| `POSTGRES_TIMEOUT`      | да          | Таймаут одного запроса к БД                 | `10s`                            |
| `POSTGRES_PORT`         |             | Порт PostgreSQL (по умолчанию `5432`)       | `5432`                           |
| `HTTP_ADDR`             | да          | Адрес HTTP-сервера                          | `:5050`                          |
| `HTTP_SHUTDOWN_TIMEOUT` |             | Таймаут graceful shutdown (по умолчанию `30s`) | `30s`                         |
| `ALLOWED_ORIGINS`       | да          | Разрешённые CORS-origins через запятую      | `http://localhost:5050,null`     |
| `LOGGER_LEVEL`          |             | `DEBUG`, `INFO`, `WARN` или `ERROR`         | `DEBUG`                          |
| `TIME_ZONE`             |             | Часовой пояс IANA (по умолчанию `UTC`)      | `Europe/Moscow`                  |

`POSTGRES_HOST`, `LOGGER_FOLDER` и `PROJECT_ROOT` задавать не нужно: их выставляют `Makefile` и `docker-compose.yaml`.

---

## Структура проекта

```
.
├── cmd/todoapp/            # main.go — точка входа и сборка зависимостей; Dockerfile
├── internal/
│   ├── core/               # общий код, не привязанный к фичам
│   │   ├── config/         # общая конфигурация (часовой пояс)
│   │   ├── domain/         # сущности Task, User, Statistics и их бизнес-правила
│   │   ├── errors/         # общие ошибки: ErrNotFound, ErrInvalidArgument, ErrConflict
│   │   ├── logger/         # zap-логгер, передаётся через context
│   │   ├── repository/postgres/pool/   # интерфейс пула БД и реализация на pgx
│   │   └── transport/http/ # сервер, роутинг, middleware, чтение запросов и запись ответов
│   └── features/           # фичи: каждая разделена на repository / service / transport
│       ├── users/
│       ├── tasks/
│       ├── statistics/
│       └── web/            # отдаёт public/index.html
├── migrations/             # SQL-миграции
├── public/index.html       # веб-интерфейс
├── docs/                   # сгенерированная Swagger-спецификация (не править вручную)
├── postman_collection.json # коллекция запросов для Postman
├── docker-compose.yaml
└── Makefile
```

## Архитектура

Каждая фича состоит из трёх слоёв. Зависимости собираются вручную в `cmd/todoapp/main.go`:

```
transport (HTTP)  →  service (бизнес-логика)  →  repository (PostgreSQL)
                            ↓
                  domain (сущности и правила)
```

Что полезно знать при работе с кодом:

- **Интерфейсы объявлены у потребителя.** Интерфейс репозитория лежит в пакете сервиса, интерфейс сервиса — в пакете транспорта. Пакет `domain` не импортирует другие пакеты проекта.
- **Ошибки → HTTP-статусы.** Слои оборачивают ошибки из `core/errors`, а `HTTPResponseHandler.ErrorResponse` переводит их в статусы: `ErrInvalidArgument` → 400, `ErrNotFound` → 404, `ErrConflict` → 409, остальные → 500.
- **Оптимистичная блокировка.** У задач и пользователей есть поле `version`. `UPDATE` проходит, только если версия не изменилась с момента чтения, иначе API отвечает 409.
- **PATCH отличает «не передано» от `null`.** Для этого используется тип `Nullable[T]`: отсутствующее поле не меняется, `null` очищает значение.
- **Правила валидации дублируются** в `domain` и в CHECK-ограничениях миграции. При изменении схемы обновляйте оба места.
- **Порядок middleware:** CORS → RequestID → Logger → Trace → Panic. Logger полагается на `request_id`, поэтому идёт после RequestID.

### Как добавить новую фичу

1. Создайте `internal/features/<name>/{repository,service,transport/http}` по образцу `tasks`.
2. При необходимости добавьте миграцию: `make migrate-create seq=<name>`.
3. Соберите зависимости и зарегистрируйте маршруты в `cmd/todoapp/main.go`.
4. Добавьте Swagger-аннотации к обработчикам и выполните `make swagger-gen`.

---

## API

Все пути указаны относительно `/api/v1`. Полное описание с примерами есть в Swagger UI.

### Пользователи

| Метод    | Путь          | Описание                                 |
|----------|---------------|------------------------------------------|
| `POST`   | `/users`      | Создать пользователя                     |
| `GET`    | `/users`      | Список (`limit`, `offset`)               |
| `GET`    | `/users/{id}` | Получить по ID                           |
| `PATCH`  | `/users/{id}` | Частично обновить                        |
| `DELETE` | `/users/{id}` | Удалить                                  |

### Задачи

| Метод    | Путь          | Описание                                         |
|----------|---------------|--------------------------------------------------|
| `POST`   | `/tasks`      | Создать задачу                                   |
| `GET`    | `/tasks`      | Список (`limit`, `offset`, фильтр `user_id`)     |
| `GET`    | `/tasks/{id}` | Получить по ID                                   |
| `PATCH`  | `/tasks/{id}` | Частично обновить (в том числе отметить выполненной) |
| `DELETE` | `/tasks/{id}` | Удалить                                          |

### Статистика

| Метод | Путь          | Описание                                                                 |
|-------|---------------|--------------------------------------------------------------------------|
| `GET` | `/statistics` | Статистика по задачам. Фильтры: `user_id`, `from`, `to` (`YYYY-MM-DD`; `to` не включается) |
