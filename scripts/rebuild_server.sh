#!/usr/bin/env bash

docker compose --env-file=".env.${ENV}" up -d --force-recreate --build --no-deps server