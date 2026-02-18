# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.25 AS build-stage

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/app

# Deploy the application binary into a lean image
FROM alpine:latest

WORKDIR /app

COPY --from=build-stage /app/server /app/server

ENV DOCKER_API_VERSION=1.41

EXPOSE 3001

RUN apk add --no-cache libc6-compat ca-certificates

ENTRYPOINT ["/app/server"]