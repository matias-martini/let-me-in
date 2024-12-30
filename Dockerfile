FROM golang:1.23

WORKDIR /app
COPY ./src .

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "main.go"]

