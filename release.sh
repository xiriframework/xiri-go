#!/usr/bin/env bash
set -euo pipefail

export GIT_SSH_COMMAND="ssh -F /workspace/xiri/.ssh/config -i /workspace/xiri/.ssh/github/id_ed25519 -o StrictHostKeyChecking=accept-new"

if [[ -n "${1:-}" ]]; then
  VERSION="$1"
else
  LAST_TAG=$(git tag -l 'v*' --sort=-v:refname | head -n1)
  if [[ -z "$LAST_TAG" ]]; then
    VERSION="v0.0.1"
  else
    MAJOR="${LAST_TAG%%.*}"
    MAJOR="${MAJOR#v}"
    MINOR="${LAST_TAG#*.}"
    MINOR="${MINOR%.*}"
    PATCH="${LAST_TAG##*.}"
    VERSION="v${MAJOR}.${MINOR}.$((PATCH + 1))"
  fi
  echo "No version specified, bumping to $VERSION"
fi

# Validate semver format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: Version must match vX.Y.Z (e.g. v0.1.0)"
  exit 1
fi

# Check for uncommitted changes
if [[ -n "$(git status --porcelain)" ]]; then
  echo "Error: Uncommitted changes. Commit or stash first."
  exit 1
fi

# Check that the bundled Claude skill is present
if [[ ! -f skills/xiri-go-expert/SKILL.md ]]; then
  echo "Error: skills/xiri-go-expert/SKILL.md missing — refuse to release."
  exit 1
fi

# Check tag doesn't already exist
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "Error: Tag $VERSION already exists."
  exit 1
fi

# CHANGELOG-Prüfung vor den Tests, damit ein Abbruch nichts angefasst hinterlässt.
if ! grep -q '^## \[Unreleased\]' CHANGELOG.md; then
  echo "Error: '## [Unreleased]' not found in CHANGELOG.md — refuse to release."
  exit 1
fi

# Leerer Abschnitt heisst: für dieses Release ist nichts dokumentiert.
if [[ -z "$(awk '/^## \[Unreleased\]/{f=1;next} /^## \[/{f=0} f' CHANGELOG.md | tr -d '[:space:]')" ]]; then
  echo "Error: '## [Unreleased]' in CHANGELOG.md is empty — nothing documented for $VERSION."
  exit 1
fi

# Run tests
echo "Running tests..."
go test ./...

# [Unreleased] in einen Versionsabschnitt drehen und ein leeres [Unreleased] darüber stehen
# lassen. Ohne das sammeln sich die Einträge mehrerer Releases dort an. Erst nach den Tests,
# damit ein Fehlschlag den Baum unberührt lässt -- die Prüfung auf uncommittete Änderungen
# weiter oben ist da schon durch, der Commit hier ist also der einzige.
echo "Rotating CHANGELOG section to $VERSION..."
awk -v ver="${VERSION#v}" '
  /^## \[Unreleased\]$/ && !done { print; print ""; print "## [" ver "]"; done = 1; next }
  { print }
' CHANGELOG.md > CHANGELOG.md.tmp && mv CHANGELOG.md.tmp CHANGELOG.md

git add CHANGELOG.md
git commit -m "chore(release): CHANGELOG-Abschnitt für $VERSION"

# Push current branch first
echo "Pushing current branch to origin..."
git push origin HEAD

# Create and push tag
echo "Creating tag $VERSION..."
git tag "$VERSION"
echo "Pushing tag to origin..."
git push origin "$VERSION"

# Create GitHub release
if command -v gh &>/dev/null; then
  echo "Creating GitHub release..."
  PREV_TAG=$(git tag -l 'v*' --sort=-v:refname | sed -n '2p')
  if [[ -n "$PREV_TAG" ]]; then
    NOTES=$(git log --pretty=format:"- %s" "$PREV_TAG..$VERSION")
  else
    NOTES=$(git log --pretty=format:"- %s" "$VERSION")
  fi
  gh release create "$VERSION" --title "$VERSION" --notes "$NOTES"
else
  echo "Warning: 'gh' CLI not found — skipping GitHub release creation."
  echo "Install: https://cli.github.com/"
fi

# Trigger Go module proxy
echo "Triggering Go module proxy..."
GOPROXY=proxy.golang.org go list -m "github.com/xiriframework/xiri-go@$VERSION"

echo ""
echo "Release $VERSION published!"
echo "Install: go get github.com/xiriframework/xiri-go@$VERSION"
