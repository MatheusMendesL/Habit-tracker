#!/bin/bash

docker run --rm \
  -v "$(pwd)/backend/services/$1:/src" \
  -w /src \
  sqlc/sqlc generate