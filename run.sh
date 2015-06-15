#!/usr/bin/env bash

mongod &
sed -i -e "s/localhost:8080/$DOMAIN:$PORT/g" /app/components/chat/display-chat.html
sed -i -e "s/localhost:8080/$DOMAIN:$PORT/g" /app/components/chat/submit-chat.html
sed -i -e "s/localhost:8080/$DOMAIN:$PORT/g" /app/components/chat/display-users.html
$PWD/chat "$@"