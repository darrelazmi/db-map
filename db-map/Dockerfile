FROM golang:1.24.3

WORKDIR /app

COPY . .

RUN go mod download

CMD ["go", "run", "main.go"]
