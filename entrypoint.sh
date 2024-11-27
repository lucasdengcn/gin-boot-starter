#!/bin/sh

set -e

cd /app

./gin-runner -e ${APP_ENV} -cfg ${APP_CFG}