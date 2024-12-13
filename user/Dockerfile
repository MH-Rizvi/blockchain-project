FROM golang:1.19-alpine

WORKDIR /app
COPY . /app

RUN go mod tidy

CMD ["go", "run", "main.go"]
