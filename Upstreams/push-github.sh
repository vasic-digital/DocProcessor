#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2026 Milos Vasic
#
# Push to vasic-digital GitHub

set -euo pipefail

BRANCH="${1:-master}"

echo "Pushing branch '$BRANCH' to GitHub (vasic-digital)..."
git push github "$BRANCH"
echo "Done."
