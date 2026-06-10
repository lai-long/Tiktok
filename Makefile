.PHONY: help docker-up docker-rebuild docker-restart docker-down test fmt run run-infra

help:
	@echo "Available targets:"
	@echo "  docker-up       - Start all services"
	@echo "  docker-rebuild  - Rebuild images and start all services"
	@echo "  docker-restart  - Restart business services (keep infrastructure running)"
	@echo "  docker-down     - Stop all services"
	@echo "  test            - Run all tests"
	@echo "  run             - Run all services locally (requires infrastructure running)"
	@echo "  run-infra       - Start infrastructure (MySQL, Redis, etcd) in Docker"

docker-up:
	cd docker && docker compose up -d

docker-rebuild:
	cd docker && docker compose up -d --build

docker-restart:
	@echo "Restarting business services: $(BIZ_SERVICES)..."
	cd docker && docker compose restart $(BIZ_SERVICES)
	@echo "Business services restarted."

docker-down:
	cd docker && docker compose down

test:
	go test -race -count=1 -coverprofile=coverage.out ./...

run-infra:
	@echo "Starting infrastructure (MySQL, Redis, etcd)..."
	@docker compose -f docker/docker-compose.yml up -d etcd mysql redis
	@echo "Infrastructure starting done."

lint:
	@echo "golangci-lint running"
	@golangci-lint run
	@echo "golangci-lint done"

run:
	@echo "Starting all services locally..."
	@trap 'kill 0; exit' SIGINT; \
	go run ./cmd/user/   & \
	go run ./cmd/video/  & \
	go run ./cmd/react/  & \
	go run ./cmd/social/ & \
	go run ./cmd/mfa/    & \
	sleep 2; \
	go run ./cmd/api/    & \
	go run ./cmd/ws/     & \
	echo "All services started. Press Ctrl+C to stop." ; \
	wait
