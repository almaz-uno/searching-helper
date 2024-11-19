#!/bin/sh -eu

cd /src/app

_term() {
  echo "Caught SIGTERM signal!"
  kill -TERM "$child" 2>/dev/null
}

trap _term TERM INT

if [ -f "/src/app/.env" ]; then
    . /src/app/.env
fi

go run . &

child=$!
wait "$child"
