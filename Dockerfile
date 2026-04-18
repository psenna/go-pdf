# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-pdf .

# Runtime stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-pdf .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./go-pdf"]