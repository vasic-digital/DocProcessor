#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2026 Milos Vasic
#
# Push to HelixDevelopment GitHub and GitLab

set -euo pipefail

BRANCH="${1:-master}"

echo "Pushing branch '$BRANCH' to HelixDevelopment remotes..."

for remote in helix-github helix-gitlab; do
    echo "  -> $remote"
    git push "$remote" "$BRANCH" 2>&1 || echo "  [WARN] Failed to push to $remote"
done

echo "Done."
