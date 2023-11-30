#!/bin/bash
APP_NAME="public-fire-c2-server"
id=$(sudo docker restart $APP_NAME)
sudo docker logs -f $id