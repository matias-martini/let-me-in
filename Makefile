build:
	docker build . -t let-me-in

up:
	docker run -it --rm -p8080:8080 let-me-in
