FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    apt-transport-https \
    ca-certificates \
    build-essential \
    gcc \
    g++ \
    python3 \
    python3-pip \
    openjdk-17-jdk \
    && mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://www.ucw.cz/isolate/debian/signing-key.asc -o /etc/apt/keyrings/isolate.asc \
    && echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/isolate.asc] http://www.ucw.cz/isolate/debian/ bookworm-isolate main" > /etc/apt/sources.list.d/isolate.list \
    && apt-get update \
    && apt-get install -y isolate \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main /app/main

WORKDIR /app

CMD ["/app/main"]

