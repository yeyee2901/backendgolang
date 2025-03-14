.PHONY: clean docker-up docker-down swagger run-local

clean:
	rm -rf ./sql/data
	rm -rf ./log/*
	rm -rf ./docker_log/*

docker-up:
	docker compose up -d --remove-orphans

docker-down:
	docker compose down

docker-build:
	docker compose build --no-cache

run-local: swagger
	go run ./cmd/api

swagger:
	swag init -g ./cmd/api/main.go --markdownFiles swagger-markdown --parseDependency true
