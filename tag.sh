#!/usr/bin/env bash

git tag $(git rev-parse --short HEAD)
git push --tags