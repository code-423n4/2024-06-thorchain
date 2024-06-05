#!/bin/bash

set -euo pipefail

# This script compares mimir ids to the develop branch to assert ids are immutable.
echo "Asserting mimir ID immutability..."

# skip for develop
if git merge-base --is-ancestor "$(git rev-parse HEAD)" "$(git rev-parse origin/develop)"; then
  echo "Skipping mimir ID lint for commit in develop ($(git rev-parse origin/develop))."
  exit 0
fi

# extract mimir ids in both current branch and develop
go run tools/mimir-ids/main.go >/tmp/mimir-ids-current
git checkout origin/develop
git checkout - -- tools scripts
go run tools/mimir-ids/main.go >/tmp/mimir-ids-develop
git checkout -

# print the diff, but do not fail the script
diff -u --color=always /tmp/mimir-ids-develop /tmp/mimir-ids-current || true

# assert that develop is a prefix of current
size=$(wc -c </tmp/mimir-ids-develop)
if ! cmp -n "$size" /tmp/mimir-ids-develop /tmp/mimir-ids-current; then
  echo "Mimir IDs are immutable. Do not remove existing IDs or insert new ones before the end. New IDs must be appended."
  if [[ $CI_MERGE_REQUEST_TITLE == *"#check-lint-warning"* ]]; then
    echo "Merge request contains #check-lint-warning."
  else
    echo 'Correct the change to mimir IDs or add "#check-lint-warning" to the PR description.'
    exit 1
  fi
fi
