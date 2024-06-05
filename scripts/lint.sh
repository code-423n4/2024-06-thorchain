#!/usr/bin/env bash
set -euo pipefail

die() {
  echo "ERR: $*"
  exit 1
}

# check docs version
version=$(cat version)
if ! grep "^  version: ${version}" openapi/openapi.yaml; then
  die "docs version (openapi/openapi.yaml) does not match version file ${version}"
fi

FILTER=(-e '^docs/' -e '.pb.go$' -e '^openapi/gen' -e '_gen.go')

if [ -n "$(git ls-files '*.go' | grep -v "${FILTER[@]}" | xargs gofumpt -l 2>/dev/null)" ]; then
  git ls-files '*.go' | grep -v "${FILTER[@]}" | xargs gofumpt -d 2>/dev/null
  die "Go formatting errors"
fi
go mod verify

./scripts/lint-handlers.bash

./scripts/lint-managers.bash

./scripts/lint-erc20s.bash

go run tools/analyze/main.go ./common/... ./constants/... ./x/... ./mimir/...

go run tools/lint-whitelist-tokens/main.go

./scripts/lint-versions.bash

./scripts/lint-mimir-ids.bash
