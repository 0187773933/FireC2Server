#!/bin/bash
APP_NAME="public-fire-c2-server"
sudo docker rm -f $APP_NAME || echo ""
id=$(sudo docker run -dit \
--name $APP_NAME \
--restart='always' \
-v $(pwd)/SAVE_FILES:/home/morphs/SAVE_FILES:rw \
--mount type=bind,source="$(pwd)"/config.yaml,target=/home/morphs/FireC2Server/config.yaml \
--network=6105-buttons-1 \
-p 5954:5954 \
$APP_NAME config.yaml)
sudo docker logs -f $id

# sudo docker network create 6105-buttons-1