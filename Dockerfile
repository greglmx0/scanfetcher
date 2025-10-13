FROM alpine:latest

RUN apk add --no-cache curl
RUN apk add chromium
RUN apk add go

ENV CGO_ENABLED=0
ENV ROD_BROWSER_PATH=/usr/bin/chromium

WORKDIR /app

RUN go install github.com/air-verse/air@v1.62.0 && mv /root/go/bin/air /usr/local/bin/

COPY go.mod go.sum ./
COPY third_party ./third_party
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml", "-d"]
