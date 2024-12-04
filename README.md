# Gin boot starter application

## Overview

for fast starting a microservice application on top of Gin framework.

## Project Structures

|Name|Description|Reference Link|
|----|-----------|--------------|
|config| configuration per env. including dev, sit, uat, prod etc| Yaml files |
|apis/controllers| apis implementation | follow OpenAPI 3.0 specification |
|apis/models| definition of request input models and response models | |
|core|define application errors||
|core/middlewares| request middleware, interceptor, global error handler etc| |
|core/logging| context logger| |
|core/otel| initializer OTEL tracer, metric, logging provider| |
|core/validators| initialize model fields validator, e.g Gender| |
|core/enums| define enums, e.g Gender| |
|infra|infra resources connectivity and management||
|infra/db| database connection, transaction | |
|infra/mq| MQ connection | e.g. Kafka, RabbitMQ |
|infra/cache| cache connection | e.g. Redis |
|jobs|schedule jobs implementation| e.g. cron job |
|messaging|event-driven integration between micro-services||
|messaging/events|definition of events|follow cloud event specification|
|messaging/consumers|consumer implementation, calling domain service methods||
|migrations| database schema changes| |
|persistence/entity| data entity mapping to database schema | |
|persistence/repository| data handling over database | using Context to manage transaction |
|server | implement routers, server, CDI | |
|services| business login implementation | using Context to manage transaction |
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

[Open API](https://github.com/swaggo/swag)

[Logging](https://github.com/rs/zerolog)

[RFC7807](https://tools.ietf.org/html/rfc7807)

[OTEL](https://opentelemetry.io/docs/languages/go/getting-started/)

[pprof](https://github.com/gin-contrib/pprof)

[JWT](<https://github.com/golang-jwt/jwt> <https://github.com/MicahParks/keyfunc>)

[ACL casbin] (<https://casbin.org/>)

## Practices

### Panic or return

Repository SHOULD return error, SHOULD NOT Panic.

Service MUST Panic if error and NEED to stop execution.

Global Recovery/Defer function SHOULD handle error and response.

### Column Default value

Always set DEFAULT value on database table column, so that can avoid can't convert from NULL to yyy error

## Reference

[zerolog](https://betterstack.com/community/guides/logging/zerolog/)

[validation](https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/)

[OTEL](https://signoz.io/blog/opentelemetry-gin/)

## Performance Reference

[DB Packages](https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/)
