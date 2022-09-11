#!/usr/bin/env bash

sudo apt update && sudo apt-get install -y curl wget

# From https://gist.github.com/lukechilds/a83e1d7127b78fef38c2914c4ececc3c
get_latest_release() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

TAG=$(get_latest_release  liboz/benchmark-japan-vps)
echo "Latest release is $TAG"
rm -f benchmark-japan-vps.tar.gz
wget https://github.com/liboz/benchmark-japan-vps/releases/download/$TAG/benchmark-japan-vps.tar.gz
tar -xf benchmark-japan-vps.tar.gz

echo "Add systemd service"
cp server/benchmark-server.service /etc/systemd/system 
systemctl daemon-reload
systemctl enable benchmark-server.service
systemctl restart benchmark-server.service
sudo reboot