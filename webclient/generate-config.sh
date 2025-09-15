#!/bin/bash

read -p "Web client port: " WEB_CLIENT_PORT
echo

cat <<EOF >webclient.json
{
  "port": $WEB_CLIENT_PORT
}
EOF

echo "Config saved to webclient.json"
