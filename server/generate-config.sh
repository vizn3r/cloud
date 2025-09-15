#!/bin/bash

read -p "Web client host: " WEB_CLIENT_HOST
read -p "Web client host port: " WEB_CLIENT_PORT
echo

cat <<EOF >server.json
{
  "webClient": {
    "host": "$WEB_CLIENT_HOST",
    "port": $WEB_CLIENT_PORT
  }
}
EOF

echo "Config saved to server.json"
