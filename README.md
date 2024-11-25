# Gin boot starter application

## Overview

for fast starting a application on top of Gin framework.

## How to run

clone this repo, and rename directory to gin001.

## Project Structures

|Name|Description|Reference Link|
|----|-----------|--------------|
|config| configuration per env. including dev, sit, uat, prod etc| Yaml files |
|apis/controllers| apis implementation | follow OpenAPI 3.0 specification |
|apis/middlewares| request middleware, interceptor etc| |
|apis/models| definition of request input models and response models | |
|infra/db| database connection | |
|infra/mq| MQ connection | e.g. Kafka, RabbitMQ |
|infra/cache| cache connection | e.g. Redis |
|jobs|schedule jobs implementation| e.g. cron job |
|messaging/events|definition of events|follow cloud event specification|
|messaging/consumers|consumer implementation, calling domain service methods||
|migrations| database schema changes| |
|persistence/entity| data entity mapping to database schema | |
|persistence/repository| data handling over database | |
|server | implement routers, server, CDI | |
|services| business login implementation | |
|wire-config| CDI initializer | |

## Schema

in General, migration files should be put in a dedicated repo, and run with a dedicated restricted User.

### migrations

[Go migrate](https://github.com/golang-migrate/migrate?tab=readme-ov-file)

## Tech stack

[PostgreSQL](https://pkg.go.dev/github.com/jackc/pgx/v5@v5.0.4/stdlib)

[SqlX](https://jmoiron.github.io/sqlx/)

[Gin](https://gin-gonic.com/docs/introduction/)

[Go migrate](https://github.com/golang-migrate/migrate)

[testify](https://github.com/stretchr/testify)

[yaml](https://github.com/go-yaml/yaml)

[Viper](https://github.com/spf13/viper)

[CDI](https://github.com/google/wire)
