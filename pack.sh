#!/usr/bin/env bash

set -e
mkdir -p bin

for GOOS in darwin linux windows; do
  for GOARCH in amd64 arm64; do
    export GOOS=$GOOS GOARCH=$GOARCH

    name="replit-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
      name="$name.exe"
    fi

    echo "Building: $name"
    xk6 build --with github.com/grafana/xk6-replit=.
    mv k6 assets
    go build -o "bin/$name" ./assets/replit.go
  done
done
