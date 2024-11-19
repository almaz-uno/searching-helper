#!/bin/sh -eu

cd /app/src

_term() {
  echo "Caught SIGTERM signal!"
  kill -TERM "$child" 2>/dev/null
}

trap _term TERM INT

if [ -f "/app/src/.env" ]; then
    . /app/src/.env
fi

go run . &

child=$!
wait "$child"
