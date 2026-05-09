# ⚡ URL Shortener

Это современный, высокопроизводительный сервис для сокращения ссылок, построенный на архитектуре Clean Architecture. Проект включает в себя полноценный API, систему аналитики в реальном времени, асинхронную обработку данных и стильный веб-интерфейс.

![Frontend Preview](web/index.html) <!-- Ссылка на файл для контекста -->

## ✨ Особенности

- **Clean Architecture**: Четкое разделение ответственности между доменом, кейсами и инфраструктурой.
- **Высокая производительность**: Использование Redis для кэширования и асинхронной записи.
- **Аналитика**: Отслеживание кликов по датам и браузерам пользователей.
- **Асинхронность**: Обработка статистики и обновление кэша через воркеры (статистика не замедляет редирект).
- **Premium UI**: Стильный фронтенд с эффектом Glassmorphism и темной темой.
- **Observability**: Интеграция с Prometheus (метрики) и OpenTelemetry (трассировка).

## 🛠 Технологический стек

- **Language**: Go 1.22+
- **Framework**: Gin Gonic
- **Database**: PostgreSQL (хранение ссылок)
- **Cache**: Redis (кэширование и асинхронная очередь)
- **Logging**: Zerolog
- **Monitoring**: Prometheus & Grafana
- **Tracing**: Jaeger (OpenTelemetry)
- **Frontend**: Vanilla JS, Modern CSS (Grid/Flexbox/Gradients)

## 🏗 Структура проекта

```text
.
├── api/                # OpenAPI спецификации
├── cmd/                # Точка входа в приложение
├── config/             # Конфигурация приложения (envconfig)
├── internal/
│   ├── app/            # Инициализация и зависимости (DI)
│   └── shortener/      # Основной домен сервиса
│       ├── controller/ # HTTP и Kafka хендлеры
│       ├── domain/     # Бизнес-логика и сущности
│       ├── dto/        # Data Transfer Objects
│       ├── usecase/    # Сценарии использования
│       └── worker/     # Фоновые задачи (Redis writer)
├── pkg/                # Переиспользуемые пакеты (logger, otel, db)
├── migrations/         # SQL миграции
├── web/                # Фронтенд файлы (HTML/CSS/JS)
└── docker-compose.yml  # Инфраструктура в Docker
```

## 🚀 Быстрый запуск

### 1. Подготовка окружения
Убедитесь, что у вас установлены Docker и Docker Compose.

### 2. Запуск инфраструктуры
```bash
docker-compose up -d
```

### 3. Запуск приложения (локально)
```bash
go run cmd/main.go
```
Сервис будет доступен по адресу: `http://localhost:8080`
Фронтенд: `http://localhost:8080/static/`

## 📡 API Endpoints

### Shorten URL
`POST /api/v1/shorten`
```json
{
  "original_url": "https://very-long-url.com",
  "shorten_code": "my-link"
}
```

### Redirect
`GET /api/v1/s/{code}`
Перенаправляет на оригинальный URL.

### Analytics
`GET /api/v1/analytics/{code}`
Возвращает статистику переходов.

## 🧪 Тестирование

Запуск всех тестов (Unit + Integration):
```bash
go test ./...
```

## 📈 Метрики и Трассировка
- **Метрики**: `http://localhost:8080/metrics`
- **Jaeger UI**: `http://localhost:8080:16686` (если запущен через docker-compose)

---
