#!/usr/bin/env bash
# Builds update.json for a Case Clicker release.
#
#   ./make-release.sh 1.3.0 "What changed, in one sentence"
#
# Produces release/update.json next to a copy of the installer, both of which
# get attached to the GitHub release. The URL it writes is GitHub's permanent
# "latest release" link, so it never needs changing between versions.

set -euo pipefail
VERSION="${1:-}"
NOTES="${2:-}"
REPO="${REPO:-YOUR-USERNAME/case-clicker}"
INSTALLER="Case Clicker Setup.exe"

if [ -z "$VERSION" ]; then
  echo "usage: ./make-release.sh <version> [notes]" >&2; exit 1
fi
if [ ! -f "$INSTALLER" ]; then
  echo "error: '$INSTALLER' not found - build it first" >&2; exit 1
fi

mkdir -p release
# GitHub rewrites spaces in asset filenames, so publish it without them.
cp "$INSTALLER" "release/CaseClickerSetup.exe"
SHA=$(sha256sum "release/CaseClickerSetup.exe" | cut -d' ' -f1)

cat > release/update.json <<JSON
{
  "version": "$VERSION",
  "notes": $(printf '%s' "$NOTES" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))'),
  "url": "https://github.com/$REPO/releases/latest/download/CaseClickerSetup.exe",
  "sha256": "$SHA",
  "date": "$(date -u +%Y-%m-%d)"
}
JSON

echo "release/update.json written:"
cat release/update.json
echo
echo "Now create a GitHub release tagged v$VERSION and attach BOTH files from release/:"
echo "  - CaseClickerSetup.exe"
echo "  - update.json"
