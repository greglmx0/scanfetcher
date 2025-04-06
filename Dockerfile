FROM golang:1.23-alpine
RUN apk add --no-cache build-base

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml", "-d"]