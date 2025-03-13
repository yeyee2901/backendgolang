.PHONY: clean docker-up docker-down

clean:
	rm -rf ./sql/data

docker-up: clean
	docker-compose up -d --remove-orphans

docker-down:
	docker-compose down
