FROM golang:1.23-bookworm

WORKDIR /app/src

COPY . .

ENTRYPOINT [ "/app/src/docker-entrypoint.sh" ]
