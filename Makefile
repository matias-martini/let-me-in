build: Dockerfile src/go.mod 
	docker compose build  
	echo "" >> build

up: build
	docker compose up 

migrate: build
	docker compose run --rm backend go run main.go db migrate

test: build
	docker compose run -e DB_NAME=test_db --rm backend go test ./... 

test-db: build
	docker compose run -e DB_NAME=test_db --rm backend go run main.go db migrate
 
fmt:
	docker compose run --rm backend go fmt ./...
