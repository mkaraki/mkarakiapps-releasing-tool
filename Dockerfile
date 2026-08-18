FROM golang:1.26-alpine AS prep

WORKDIR /app

#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

COPY . .

FROM prep AS builder-pr-version-validator

RUN go build -o pr-version-validator cmd/pr-version-validate/main.go

FROM alpine

WORKDIR /workdir

COPY --from=builder-pr-version-validator /app/pr-version-validator /usr/bin/pr-version-validator

COPY entrypoint-*.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint-*.sh
