#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2026 Milos Vasic
#
# Push to all 4 remotes: vasic-digital GitHub/GitLab + HelixDevelopment GitHub/GitLab

set -euo pipefail

BRANCH="${1:-master}"

echo "Pushing branch '$BRANCH' to all remotes..."

for remote in github gitlab helix-github helix-gitlab; do
    echo "  -> $remote"
    git push "$remote" "$BRANCH" 2>&1 || echo "  [WARN] Failed to push to $remote"
done

echo "Done."
