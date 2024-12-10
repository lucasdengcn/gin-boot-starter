#!/bin/sh

set -e

cd /app

./app-runner -e ${APP_ENV} -w ${APP_BASE}