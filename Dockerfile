# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-pdf .

# Runtime stage
FROM debian:bookworm-slim AS runtime
WORKDIR /app
COPY --from=builder /app/go-pdf /app/go-pdf
COPY --from=builder /app/templates /app/templates
EXPOSE 8080
CMD ["/app/go-pdf"]