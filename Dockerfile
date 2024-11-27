ARG GO_VERSION=1.23.3

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /build/gin001
WORKDIR /build/gin001

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

COPY --from=builder /build/gin001/gin-runner .
COPY --from=builder /build/gin001/entrypoint.sh .
COPY --from=builder /build/gin001/config/*.yaml ./config/
COPY --from=builder /build/gin001/migrations ./migrations/

EXPOSE 8080

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]