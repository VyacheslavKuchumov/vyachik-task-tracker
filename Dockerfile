FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/server ./cmd/main.go && \
    CGO_ENABLED=0 go build -o /out/app-migrate ./cmd/migrate/main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/server /usr/local/bin/server
COPY --from=builder /out/app-migrate /usr/local/bin/app-migrate
COPY --from=builder /src/cmd/migrate/migrations ./cmd/migrate/migrations

EXPOSE 8000

CMD ["sh", "-c", "app-migrate up && server"]
