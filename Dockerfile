# syntax=docker/dockerfile:1

# Step 1: Build stage
FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    tesseract-ocr-tha \
    libtesseract-dev \
    libleptonica-dev \
    imagemagick \
    gcc g++ \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Module caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source
COPY . .

WORKDIR /app/cmd
RUN go build -o bot .

# Step 2: Runtime stage (slim)
FROM debian:bullseye-slim AS runtime

RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    tesseract-ocr-tha \
    imagemagick \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/cmd/bot .

# Uncomment the following line to use .env.deploy (local build only!)
# COPY .env.deploy ./

CMD ["./bot"]
