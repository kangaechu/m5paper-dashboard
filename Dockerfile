FROM docker.io/library/busybox:1.38-uclibc AS busybox

FROM docker.io/rclone/rclone:1 AS rclone

FROM golang:1.26.4-trixie AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o m5paper-dashboard ./cmd/local

FROM gcr.io/distroless/base-debian13

COPY --from=busybox /bin/busybox /bin/busybox
RUN ["/bin/busybox", "--install", "/bin"]

USER nonroot:nonroot
WORKDIR /app
COPY --chown=nonroot:nonroot --from=builder /build/m5paper-dashboard /app/m5paper-dashboard
COPY --chown=nonroot:nonroot --from=rclone /usr/local/bin/rclone /usr/local/bin/rclone
COPY --chown=nonroot:nonroot run.sh /app/run.sh

# ENV SYNC_CLOUD_STORAGE="true"
# ENV RCLONE_DESTINATION
# ENV RCLONE_CONFIG

ENTRYPOINT ["/bin/sh", "-c", "/app/run.sh"]
