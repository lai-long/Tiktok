#!/bin/bash

# Build API service
mkdir -p output/bin
cp script/* output 2>/dev/null
chmod +x output/bootstrap.sh
go build -o output/bin/api ./cmd/api/
go build -o output/bin/ws ./cmd/ws/

echo "Build completed: api, ws"