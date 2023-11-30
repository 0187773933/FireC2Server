#!/bin/bash

# so you have to run it once , without mounting the ADB_KEYS folder
# run some command , and pair with the fire cube
# then pull them locally :
	# sudo docker cp 05099c8c4137:/home/morphs/.android ADB_KEYS

APP_NAME="public-fire-c2-server"
sudo docker rm -f $APP_NAME || echo ""
id=$(sudo docker run -dit \
--name $APP_NAME \
--restart='always' \
-v $(pwd)/SAVE_FILES:/home/morphs/SAVE_FILES:rw \
-v $(pwd)/ADB_KEYS:/home/morphs/.android:r \
--mount type=bind,source="$(pwd)"/config.yaml,target=/home/morphs/config.yaml \
--network=6105-buttons-1 \
-p 5954:5954 \
$APP_NAME /home/morphs/config.yaml)
sudo docker logs -f $id

# sudo docker network create 6105-buttons-1