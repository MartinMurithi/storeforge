#!/bin/sh
# -----------------------------------------------------------------------------
# Service test runner
# -----------------------------------------------------------------------------
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SERVICE_NAME=$(basename $(pwd))
echo -e "${YELLOW}--------------------------------------------------${NC}"
echo -e "${YELLOW}Running tests for service: $SERVICE_NAME${NC}"
echo -e "${YELLOW}--------------------------------------------------${NC}"

# Find all Go test packages in this service
PACKAGES=$(go list ./...)

if [ -z "$PACKAGES" ]; then
  echo -e "${YELLOW}No Go packages found in $SERVICE_NAME, skipping tests.${NC}"
  exit 0
fi

# Run tests
for pkg in $PACKAGES; do
  echo -e "${YELLOW}Testing package: $pkg${NC}"
  go test -v -race "$pkg"
done

echo -e "${GREEN}All tests passed for $SERVICE_NAME!${NC}"

