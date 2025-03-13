.PHONY: clean docker-up docker-down swagger run-local

clean:
	rm -rf ./sql/data

docker-up:
	docker-compose up -d --remove-orphans

docker-down:
	docker-compose down

run-local: swagger
	go run ./cmd/api

swagger:
	swag init -g ./cmd/api/main.go --markdownFiles swagger-markdown --parseDependency true
