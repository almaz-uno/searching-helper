#!/bin/bash -eu

cd $(dirname $(realpath $0))/..

APP_NAME=$(basename $(realpath .))

TAG=$(git describe --tags --abbrev=0)

TARGET_IMG=mkovrov/$APP_NAME

docker build --file Dockerfile -t $TARGET_IMG:$TAG -t $TARGET_IMG:latest .

docker login

docker push $TARGET_IMG:$TAG
docker push $TARGET_IMG:latest
