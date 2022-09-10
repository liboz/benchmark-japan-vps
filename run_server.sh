#!/usr/bin/env bash

# From https://gist.github.com/lukechilds/a83e1d7127b78fef38c2914c4ececc3c
get_latest_release() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

DEBUG_FLAG=${1:-false}
TAG=$(get_latest_release  liboz/benchmark-japan-vps)
echo "Latest release is $TAG"
wget https://github.com/liboz/benchmark-japan-vps/releases/download/$TAG/benchmark-japan-vps.tar.gz
tar -xf benchmark-japan-vps.tar.gz
./benchmark-japan-vps $DEBUG_FLAG
