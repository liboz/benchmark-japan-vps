#!/usr/bin/env bash
set -e
# reads from tools.go and installs all tools
cat tools.go | grep _ | awk -F'"' '{print $2}' | xargs -tI % go install %