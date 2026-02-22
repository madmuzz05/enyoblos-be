# ===== BUILD STAGE =====
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

# ===== FINAL STAGE =====
FROM alpine:latest

WORKDIR /app

# install cert
RUN apk add --no-cache ca-certificates

# 🔐 buat user non-root
RUN adduser -D appuser

# copy binary
COPY --from=builder /app/app .

# copy migrations
COPY --from=builder /app/package/database/postgres/migrations ./package/database/postgres/migrations

# kasih permission ke user baru
RUN chown -R appuser:appuser /app

# 🔐 pakai user non-root
USER appuser

CMD ["./app"]