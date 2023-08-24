#!/usr/bin/env bash
# 2021-07-08 WATERMARK, DO NOT REMOVE - This script was generated from the Kurtosis Bash script template

set -euo pipefail   # Bash "strict mode"
script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
server_root_dirpath="$(dirname "${script_dirpath}")"



# ==================================================================================================
#                                             Constants
# ==================================================================================================
BUILD_DIRNAME="build"

MAIN_GO_FILEPATH="${server_root_dirpath}/main.go"
MAIN_BINARY_OUTPUT_FILENAME="kurtosis-package-indexer"
MAIN_BINARY_OUTPUT_FILEPATH="${server_root_dirpath}/${BUILD_DIRNAME}/${MAIN_BINARY_OUTPUT_FILENAME}"

# ==================================================================================================
#                                             Main Logic
# ==================================================================================================

# Test code
echo "Running unit tests..."
if ! cd "${server_root_dirpath}"; then
  echo "Couldn't cd to the server root dirpath '${server_root_dirpath}'" >&2
  exit 1
fi
if ! CGO_ENABLED=0 go test "./..."; then
  echo "Tests failed!" >&2
  exit 1
fi
echo "Tests succeeded"

# Build binary for packaging inside an Alpine Linux image
echo "Building server main.go '${MAIN_GO_FILEPATH}'..."
if ! CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${MAIN_BINARY_OUTPUT_FILEPATH}" "${MAIN_GO_FILEPATH}"; then
  echo "Error: An error occurred building the server code" >&2
  exit 1
fi
echo "Successfully built server code"
