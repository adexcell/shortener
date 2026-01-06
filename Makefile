DC := docker compose
PROJECT := shortener

.PHONY:
migrate-up: ## Запуск migrate up
	migrate -path ./migrations \
	-database "postgres://postgres:pass@localhost:5432/shortener?sslmode=disable" up

.PHONY:
migrate-force-0: ## Запуск migrate force 0
	migrate -path ./migrations \
	-database "postgres://postgres:pass@localhost:5432/shortener?sslmode=disable" force 0

.PHONY:
migrate-force-1: ## Запуск migrate force 0
	migrate -path ./migrations \
	-database "postgres://postgres:pass@localhost:5432/shortener?sslmode=disable" force 1


.PHONY:
up:      ## Запуск в фоне 	
	$(DC) -p $(PROJECT) up -d 

.PHONY:	
down:    ## Полная остановка (контейнеры + сети) 	
	$(DC) -p $(PROJECT) down 

.PHONY:	
stop:    ## Остановка без удаления 	
	$(DC) -p $(PROJECT) 

.PHONY:
stop-start:   ## Запуск остановленных контейнеров 	
	$(DC) -p $(PROJECT) start

.PHONY:
logs:    ## Логи всех сервисов (следить) 	
	$(DC) logs -f 

.PHONY:	
logs-web: ## Логи конкретного сервиса 	
	$(DC) logs -f web 

.PHONY:	
ps:      ## Статус контейнеров 	
	$(DC) ps 

.PHONY:
exec-web: ## Войти в контейнер web 	
	$(DC) exec web bash

.PHONY:
clean:   ## Удаление volumes + images 	
	$(DC) down -v && docker system prune -f 

.PHONY:
run: ## Запуск main.go
	go run ./cmd/shortener/main.go

.PHONY:
help: ## Список команд
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
