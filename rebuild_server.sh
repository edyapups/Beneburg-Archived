#!/usr/bin/env bash

export $(grep -v '^#' .env | xargs)

docker-compose up -d --force-recreate --build --no-deps server