#!/bin/sh
set -e

echo "======================================"
echo " Running Go test suite"
echo "======================================"

go test ./...

echo "======================================"
echo " Test suite passed successfully"
echo "======================================"
