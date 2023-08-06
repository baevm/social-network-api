#!/bin/sh

source ".env"

set -e

echo "running db migration"
./migrate -path ./migrations -database "$DB_DSN" -verbose up

echo "starting app"
exec "$@"