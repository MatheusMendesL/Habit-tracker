#!/bin/bash

MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$(pwd)/backend/services/$1:/src" \
  -w /src \
  sqlc/sqlc generate
