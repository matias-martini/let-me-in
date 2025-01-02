build: Dockerfile src/go.mod 
	docker compose build  
	echo "" >> build

up: build
	docker compose up 

migrate: build
	docker compose run --rm backend go run main.go db migrate

 
