#!/usr/bin/env bash
set -euo pipefail

# prompt for gitlab merge request id
read -rp "Enter Gitlab Merge Request ID: " MR

BRANCH="mr-$MR"

git branch -D "${BRANCH}" || true
git fetch origin merge-requests/"$MR"/head:"${BRANCH}"
git checkout "${BRANCH}"
git push --set-upstream origin "${BRANCH}" -f --no-verify
git checkout "@{-1}"

echo
echo "Navigate to https://gitlab.com/thorchain/thornode/-/pipelines/new and run the pipeline for branch: mr-$MR"
