#!/bin/sh

# DB env vars in the format the pgdb image wants
export POSTGRES_USER=foo
export POSTGRES_PASSWORD=verysecret
export POSTGRES_DB=bar

# DB env vars in the format the migrations and apps want
export PGHOST=localhost
export PGPORT=5432
export PGUSER=foo
export PGPASSWORD=verysecret
export PGDATABASE=bar

# Extra app env vars
export SERVE_IP="0.0.0.0"
export SERVE_PORT="8000"
export SERVE_WRITE_TIMEOUT="10s"
export SERVE_READ_TIMEOUT="10s"

