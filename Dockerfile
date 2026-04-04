FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o m5paper-dashboard ./cmd/local

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /build/m5paper-dashboard /app/m5paper-dashboard
ENTRYPOINT ["/app/m5paper-dashboard"]
