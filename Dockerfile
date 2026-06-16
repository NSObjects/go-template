ARG GO_VERSION=1.26.1

FROM golang:${GO_VERSION}-alpine AS builder

# ca-certificates is required to call HTTPS endpoints.
# tzdata is required for time zone info.
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

WORKDIR /src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -trimpath -ldflags="-w -s" -o /out/app ./main.go

FROM scratch AS final

ENV LANG=en_US.UTF-8

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/app/configs/config.toml /configs/config.toml
COPY --from=builder /out/app /app

EXPOSE 9322

ENTRYPOINT ["/app", "--config", "/configs/config.toml"]
