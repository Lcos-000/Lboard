.PHONY: dev-up dev-down dev-logs server-run server-test server-tidy web-install web-dev web-build check
dev-up:
	docker compose up -d --build
dev-down:
	docker compose down
dev-logs:
	docker compose logs -f
server-tidy:
	cd server && go mod tidy
server-run:
	cd server && go run ./cmd/whiteboard-server
server-test:
	cd server && go test ./...
web-install:
	cd web && pnpm install
web-dev:
	cd web && pnpm dev
web-build:
	cd web && pnpm build
check:
	curl -s http://localhost:8080/healthz
	@echo ""
	curl -s http://localhost:8080/readyz
	@echo ""