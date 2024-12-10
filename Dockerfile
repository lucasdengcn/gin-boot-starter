ARG GO_VERSION=1.23.3

FROM golang:${GO_VERSION}-alpine AS builder

ARG APP_NAME=gin-boot-starter

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

# match with mod's name
RUN mkdir -p /build/${APP_NAME}
WORKDIR /build/${APP_NAME}

ENV GIN_MODE=release

COPY go.mod go.sum entrypoint.sh ./
RUN dos2unix entrypoint.sh
RUN go mod download

#
COPY api ./api
COPY core ./core 
COPY docs ./docs
COPY infra ./infra 
COPY job ./job
COPY messaging ./messaging
COPY persistence ./persistence
COPY server ./server
COPY service ./service
COPY wire-config ./wire-config
COPY main.go ./

RUN go build -o app-runner .

#
COPY migration ./migration
COPY config ./config 
COPY secrets ./secrets


FROM alpine:latest

ARG APP_NAME=gin-boot-starter

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
RUN mkdir -p /app/config
RUN mkdir -p /app/secrets

ENV APP_ENV dev
ENV APP_BASE /app
ENV APP_PORT 8080

WORKDIR /app

COPY --from=builder /build/${APP_NAME}/app-runner .
COPY --from=builder /build/${APP_NAME}/entrypoint.sh .
COPY --from=builder /build/${APP_NAME}/config/*.yaml ./config/
COPY --from=builder /build/${APP_NAME}/secrets/*.pem ./secrets/
COPY --from=builder /build/${APP_NAME}/config/*.conf ./config/
COPY --from=builder /build/${APP_NAME}/migration ./migration/

EXPOSE ${APP_PORT}

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]