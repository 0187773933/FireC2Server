#!/bin/bash

sudo chown -R morphs:morphs /home/morphs/SAVE_FILES
sudo chown -R morphs:morphs /home/morphs/.android

HASH_FILE="/home/morphs/git.hash"
REMOTE_HASH=$(git ls-remote https://github.com/0187773933/FireC2Server.git HEAD | awk '{print $1}')
if [ -f "$HASH_FILE" ]; then
    STORED_HASH=$(sudo cat "$HASH_FILE")
else
    STORED_HASH=""
fi
if [ "$REMOTE_HASH" == "$STORED_HASH" ]; then
        echo "No New Updates Available"
        cd /home/morphs/FireC2Server
        LOG_LEVEL=debug exec /home/morphs/FireC2Server/server "$@"
else
        echo "New updates available. Updating and Rebuilding Go Module"
        echo "$REMOTE_HASH" | sudo tee "$HASH_FILE"
        cd /home/morphs
        sudo rm -rf /home/morphs/FireC2Server
        git clone "https://github.com/0187773933/FireC2Server.git"
        sudo chown -R morphs:morphs /home/morphs/FireC2Server
        cd /home/morphs/FireC2Server
        /usr/local/go/bin/go mod tidy
        GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o /home/morphs/FireC2Server/server
        LOG_LEVEL=debug exec /home/morphs/FireC2Server/server "$@"
fi