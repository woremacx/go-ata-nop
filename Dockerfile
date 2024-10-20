FROM woremacx/golang:1.23 AS builder

COPY ./ /build
WORKDIR /build

# RUN go mod tidy -v

RUN go build -v


# Final stage
FROM scratch

# Copy the Go binary
COPY --from=builder /build/go-ata-nop /app/go-ata-nop

USER root:root

CMD ["/app/go-ata-nop"]

LABEL org.opencontainers.image.source=https://github.com/woremacx/go-ata-nop
