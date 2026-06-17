ARG GO_VERSION=1.26.1

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

# ca-certificates is required to call HTTPS endpoints.
# tzdata is required for time zone info.
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
RUN addgroup -S -g 10001 app && adduser -S -D -H -u 10001 -G app app

WORKDIR /src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0

RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-$(go env GOARCH)} go build -trimpath -ldflags="-w -s" -o /out/app ./main.go

FROM scratch AS final

ENV LANG=en_US.UTF-8
ENV GO_TEMPLATE_SYSTEM_LEVEL=2 \
    GO_TEMPLATE_LOG_FORMAT=json

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /src/app/configs/config.toml /configs/config.toml
COPY --from=builder /out/app /app

EXPOSE 9322
USER app:app

ENTRYPOINT ["/app", "--config", "/configs/config.toml"]
