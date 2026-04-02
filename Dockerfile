FROM golang:1.26-alpine AS base

WORKDIR /app

RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# =========================
# 🔥 DEV (Hot Reload)
# =========================
FROM base AS dev

RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 8080

CMD ["air"]

# =========================
# 🏗️ BUILD (produção)
# =========================
FROM base AS builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o server \
    ./cmd/payment_processor

# ---- Runtime stage ----
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["./server"]
