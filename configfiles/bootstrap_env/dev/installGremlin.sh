#!/bin/bash

#
# Gremlin installation; see instructions as per:
# https://confluence.mathworks.com/display/WS/Gremlin+Set+Up+and+Test+Run+via+Gremlin+UI
#

echo "deb https://deb.gremlin.com/ release non-free" | sudo tee /etc/apt/sources.list.d/gremlin.list
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 9CDB294B29A5B1E2E00C24C022E8EF3461A50EF6
sudo apt-get update && sudo apt-get install -y gremlin gremlind

#
# Configuration
#
sudo cp /tmp/bootstrap/files/config.yaml /etc/gremlin/config.yaml

sudo systemctl start gremlind.service
sudo systemctl restart gremlind