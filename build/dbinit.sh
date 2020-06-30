#!/bin/sh

set -e

# Spin until pg is up and started
until pg_isready; do
  echo >&2 "pg not ready"
  sleep 1
done

# Run the db migrations
echo >&2 "running migrations"
tern migrate -m migrations

# Run the db seed application
echo >&2 "running db seed"
seed
