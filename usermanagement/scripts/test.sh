#!/bin/sh
set -e

echo "[usermanagement] running tests..."

go test ./...

echo "[usermanagement] tests passed"
