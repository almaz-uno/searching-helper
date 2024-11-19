#!/bin/bash -eu

cd $(dirname $(realpath $0))/..

APP_NAME=$(basename $(realpath .))

TARGET_IMG=mkovrov/$APP_NAME

docker pull $TARGET_IMG

docker rm -f $APP_NAME

docker run -d --restart=unless-stopped \
    -p 32080:32080 --name $APP_NAME \
    -v /root/.acme.sh:/root/.acme.sh \
    -e LISTEN_ON=:32080 \
    -e CERT_FILE=/root/.acme.sh/ilovlya.space_ecc/fullchain.cer \
    -e KEY_FILE=/root/.acme.sh/ilovlya.space_ecc/ilovlya.space.key \
    docker.io/mkovrov/searching-helper \
