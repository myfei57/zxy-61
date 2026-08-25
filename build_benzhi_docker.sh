#!/usr/bin/env sh
set -eu
image="${1:-coldstore-benzhi}"
docker build -f benzhi.Dockerfile -t "$image" .
printf 'built image %s\n' "$image"
