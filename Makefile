build:
	docker build . -t let-me-in

up:
	docker compose up 

migrate:
	docker compose run backend go run main.go db migrate
