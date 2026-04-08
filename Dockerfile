# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/
COPY public/ public/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM scratch
COPY --from=builder /app/server /
COPY --from=builder /app/public/index.html /public/index.html

EXPOSE 8080
CMD ["/server"]
