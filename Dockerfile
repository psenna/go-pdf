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
# Install pdfcpu for PDF optimization
RUN apt-get update && apt-get install -y wget && \
    wget -q https://github.com/pdfcpu/pdfcpu/releases/latest/download/pdfcpu_linux_amd64.tar.gz -O - | tar -xzf - -C /usr/local/bin && \
    rm -f pdfcpu_linux_amd64.tar.gz && \
    apt-get purge -y wget && \
    apt-get clean
COPY --from=builder /app/go-pdf /app/go-pdf
COPY --from=builder /app/templates /app/templates
EXPOSE 8080
CMD ["/app/go-pdf"]