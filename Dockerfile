FROM golang:1.27-alpine AS prep

WORKDIR /app

#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

COPY . .

FROM prep AS builder-pr-version-validator

RUN go build -o pr-version-validator cmd/pr-version-validate/main.go

FROM prep AS builder-auto-tag

RUN go build -o auto-tag cmd/auto-tag/main.go

FROM alpine

RUN apk --no-cache add ca-certificates git

WORKDIR /workdir

COPY --from=builder-pr-version-validator /app/pr-version-validator /usr/bin/pr-version-validator
COPY --from=builder-auto-tag /app/auto-tag /usr/bin/auto-tag

COPY entrypoint-*.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint-*.sh
