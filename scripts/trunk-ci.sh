#!/usr/bin/env bash

# This script wraps execution of trunk when run in CI.

set -euo pipefail

FLAGS="-j8 --ci"

if [ -n "${CI_MERGE_REQUEST_ID-}" ]; then
  # if go modules or trunk settings changed, also run with --all on merge requests
  if ! git diff --exit-code origin/develop -- go.mod go.sum .trunk >/dev/null; then
    FLAGS="$FLAGS --all"
  else
    FLAGS="$FLAGS --upstream origin/develop"
  fi
else
  FLAGS="$FLAGS --all"
fi

# get directory of this script
WD=$(dirname "$0")

# run trunk
echo "Running: $WD/trunk check $FLAGS"
# trunk-ignore(shellcheck/SC2086): expanding $FLAGS as flags
"$WD"/trunk check $FLAGS
