# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /health-checker ./main.go

# Runtime stage
FROM alpine:3.21

RUN apk --no-cache add ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /health-checker /health-checker
COPY config.yaml /config.yaml

USER appuser

EXPOSE 9090

ENTRYPOINT ["/health-checker"]
