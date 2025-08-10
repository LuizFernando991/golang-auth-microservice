FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /auth ./cmd/server


FROM alpine:3.18
RUN apk add --no-cache ca-certificates postgresql-client bash

WORKDIR /app

COPY --from=builder /auth /auth
COPY internal/db/migrations.sql ./internal/db/migrations.sql
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/auth"]
