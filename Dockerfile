# Build stage
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o backend.exe .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend.exe .
COPY --from=builder /app/static ./static

# Expose port and set entrypoint
EXPOSE 8080
CMD ["./backend.exe"]