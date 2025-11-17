#!/usr/bin/env bash
set -euo pipefail

export GOOS=${GOOS:-linux}
export GOARCH=${GOARCH:-arm64}
export CGO_ENABLED=0
export GOPROXY=${GOPROXY:-off}
export GOSUMDB=${GOSUMDB:-off}

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "${ROOT_DIR}"

echo "Building todo API binary"
mkdir -p bin
go build -o bin/todo-app ./cmd/server

echo "Building integration client binary"
mkdir -p integration/bin
go build -o integration/bin/todo-integration ./integration

echo "Binaries ready under bin/ and integration/bin"
