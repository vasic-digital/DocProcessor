#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2026 Milos Vasic
#
# Push to vasic-digital GitLab

set -euo pipefail

BRANCH="${1:-master}"

echo "Pushing branch '$BRANCH' to GitLab (vasic-digital)..."
git push gitlab "$BRANCH"
echo "Done."
