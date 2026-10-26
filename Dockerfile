# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/local-offers-api ./cmd/api

FROM alpine:3.22

RUN apk add --no-cache ca-certificates wget \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /data/uploads \
    && chown -R app:app /data/uploads

WORKDIR /app
COPY --from=build /out/local-offers-api /app/local-offers-api

USER app
EXPOSE 8080

ENTRYPOINT ["/app/local-offers-api"]
