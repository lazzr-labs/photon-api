# !/bin/bash

docker buildx build --file Dockerfile . --tag ghcr.io/lazzr-labs/photon:ute1783536597 --platform linux/amd64,linux/arm64 --push
