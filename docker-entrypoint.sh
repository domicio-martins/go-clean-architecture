#!/bin/sh
set -e

if [[ ! -z "${APP_STAGE}" ]]; then
    eval $(curl -s env.getter/${APP_NAME}?format=bash)
fi

exec /app/go-clean-architecture