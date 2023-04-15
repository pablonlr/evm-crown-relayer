FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY config.json .

RUN go build -o main ./cmd/main/main.go

# Ejecutar la aplicación
CMD ["./main"]
