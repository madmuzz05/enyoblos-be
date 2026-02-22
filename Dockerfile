# ---------- BUILD STAGE ----------
FROM golang:alpine AS builder

WORKDIR /app

# install git (biar go mod aman)
RUN apk add --no-cache git

# copy dependency dulu (biar cache optimal)
COPY go.mod go.sum ./
RUN go mod download

# copy semua source code
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# ---------- FINAL STAGE ----------
FROM alpine:latest

WORKDIR /app

# install timezone (optional tapi bagus)
RUN apk add --no-cache tzdata

# copy binary
COPY --from=builder /app/app .

# copy folder penting (INI YANG FIX BUG KAMU 🔥)
COPY --from=builder /app/package ./package

# expose port
EXPOSE 8080

# run app
CMD ["./app"]