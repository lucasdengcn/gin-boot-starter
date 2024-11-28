ARG GO_VERSION=1.23.3

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /build/gin-boot-starter
WORKDIR /build/gin-boot-starter

ENV GIN_MODE=release

COPY go.mod go.sum main.go entrypoint.sh ./
#
COPY apis ./apis
COPY config ./config 
COPY core ./core 
COPY docs ./docs
COPY infra ./infra 
COPY jobs ./jobs
COPY messaging ./messaging
COPY migrations ./migrations
COPY persistence ./persistence
COPY server ./server
COPY services ./services
COPY wire-config ./wire-config

RUN dos2unix entrypoint.sh
RUN go mod download
RUN go build -o gin-runner .

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
RUN mkdir -p /app/config

ENV APP_ENV dev
ENV APP_CFG /app/config

WORKDIR /app

COPY --from=builder /build/gin-boot-starter/gin-runner .
COPY --from=builder /build/gin-boot-starter/entrypoint.sh .
COPY --from=builder /build/gin-boot-starter/config/*.yaml ./config/
COPY --from=builder /build/gin-boot-starter/config/*.pem ./config/
COPY --from=builder /build/gin-boot-starter/migrations ./migrations/

EXPOSE 8080

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]