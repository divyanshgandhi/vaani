.PHONY: test docker dev-up

test:
	go test ./...

docker:
	@echo "Building multi-arch Docker images"
	@for service in $(shell ls -d services/*/); do \
		echo "Building $$service"; \
		docker build -t vaani/$$(basename $$service) $$service; \
	done

dev-up:
	docker-compose -f infra/docker/docker-compose.yml up -d

dev-down:
	docker-compose -f infra/docker/docker-compose.yml down 