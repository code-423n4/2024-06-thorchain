#!/bin/bash

set -euo pipefail

# This script compares versioned functions to the develop branch to ensure no logic
# changes in historical versions that would cause consensus failure.

VERSION=$(awk -F. '{ print $2 }' version)
CI_MERGE_REQUEST_TITLE=${CI_MERGE_REQUEST_TITLE-}

# skip for develop
if git merge-base --is-ancestor "$(git rev-parse HEAD)" "$(git rev-parse origin/develop)"; then
  echo "Skipping lint for commit in develop ($(git rev-parse origin/develop))."
  exit 0
fi

FAILED=false

echo -n "Checking unversioned tests for versioned managers... "
go run tools/current-test-versions/main.go --managers >/tmp/testing-old-versions
if [ -s /tmp/testing-old-versions ]; then
  echo "Detected use of versioned manager in unversioned test."
  if [[ $CI_MERGE_REQUEST_TITLE == *"#check-lint-warning"* ]]; then
    echo "Merge request is marked unsafe."
  else
    echo 'In the following locations, switch to the current manager, or add "#check-lint-warning" to the PR description.'
    cat /tmp/testing-old-versions
    FAILED=true
  fi
fi
echo "OK"

go run tools/versioned-functions/main.go --version="$VERSION" >/tmp/versioned-functions-current
go run tools/versioned-tokenlists/main.go --version="$VERSION" >/tmp/versioned-tokenlists-current
git checkout origin/develop
git checkout - -- tools scripts
go run tools/versioned-functions/main.go --version="$VERSION" >/tmp/versioned-functions-develop
go run tools/versioned-tokenlists/main.go --version="$VERSION" >/tmp/versioned-tokenlists-develop
git checkout -

gofumpt -w /tmp/versioned-functions-develop /tmp/versioned-functions-current

echo "Linting versioned functions..."
if ! diff -u -F '^func' -I '^//' --color=always /tmp/versioned-functions-develop /tmp/versioned-functions-current; then
  echo "Detected change in versioned function."
  if [[ $CI_MERGE_REQUEST_TITLE == *"#check-lint-warning"* ]]; then
    echo "Merge request is marked unsafe."
  else
    echo 'Correct the change, add a new versioned function, or add "#check-lint-warning" to the PR description.'
    FAILED=true
  fi
fi

echo "Linting versioned tokenlists..."
if ! diff -u -F '^Check' --color=always /tmp/versioned-tokenlists-develop /tmp/versioned-tokenlists-current; then
  echo "Detected change in versioned tokenlist."
  FAILED=true
fi

if $FAILED; then
  echo "Lint failed."
  exit 1
fi
