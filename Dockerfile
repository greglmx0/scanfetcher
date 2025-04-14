FROM golang:1.24

RUN apt-get update && apt-get install -y \
    git \
    curl \
    bash \
    ca-certificates \
    libnss3 \
    libxss1 \
    libasound2 \
    libxtst6 \
    libgtk-3-0 \
    libgbm1 \
    libgconf-2-4 \
    libglib2.0-0 \
    xvfb \
    && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1
ENV ROD_BROWSER_PATH=/usr/bin/chromium

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml", "-d"]
