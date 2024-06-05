#!/bin/bash

set -euo pipefail

CHANGES=$(git status --porcelain)

if [ -n "$CHANGES" ]; then
  echo "Detected changes in generated code:"
  echo "$CHANGES"
  echo "Please run 'make generate' and commit the changes."
  exit 1
fi

echo "No changes detected in generated code."
