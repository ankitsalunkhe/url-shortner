.PHONY: run
run:
	go run cmd/main.go

.PHONY: infra
infra:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: generate
generate:
	go generate ./...