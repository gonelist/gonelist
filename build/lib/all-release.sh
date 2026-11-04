#!/usr/bin/env bash
set -euo pipefail

GONELIST_ROOT=${GONELIST_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)}
cd "$GONELIST_ROOT"
OUTPUT=${OUTPUT:-"$GONELIST_ROOT/release"}
DIST_VERSION=${DIST_VERSION:-v0.6.0}
RELEASE_VERSION=${RELEASE_VERSION:-$(git describe --tags --exact-match)}
COMMIT_ID=$(git rev-parse HEAD)
BUILD_DATE=${BUILD_DATE:-$(date -u +'%Y-%m-%dT%H:%M:%SZ')}
GIT_TREE_STATE=clean
if [ -n "$(git status --porcelain --untracked-files=no)" ]; then GIT_TREE_STATE=dirty; fi
mkdir -p "$OUTPUT"

if [ ! -d "$OUTPUT/dist" ]; then
  curl -fL --retry 3 "https://github.com/gonelist/gonelist-web/releases/download/${DIST_VERSION}/dist.tar.gz" -o "$OUTPUT/dist.tar.gz"
  tar -xzf "$OUTPUT/dist.tar.gz" -C "$OUTPUT"
fi

targets=(linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64 darwin/amd64 darwin/arm64)
for target in "${targets[@]}"; do
  target_os=${target%/*}
  target_arch=${target#*/}
  name="gonelist_${target_os}_${target_arch}"
  stage="$OUTPUT/$name"
  binary="$name"
  if [ "$target_os" = windows ]; then binary="$binary.exe"; fi
  mkdir -p "$stage"
  CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" go build -trimpath -o "$stage/$binary" \
    -ldflags "-s -w -X main.Version=$RELEASE_VERSION -X main.buildDate=$BUILD_DATE -X main.gitCommit=$COMMIT_ID -X main.gitTreeState=$GIT_TREE_STATE" .
  cp -a "$OUTPUT/dist" "$stage/"
  cp config.yml LICENSE SECURITY.md "$stage/"
  mkdir -p "$stage/docs/releases"
  cp "docs/releases/$RELEASE_VERSION.md" "$stage/docs/releases/"
  cp "docs/releases/$RELEASE_VERSION.md" "$stage/RELEASE_NOTES.md"
  tar -czf "$OUTPUT/$name.tar.gz" -C "$OUTPUT" "$name"
  echo "Built $target"
done
(cd "$OUTPUT" && sha256sum gonelist_*.tar.gz > SHA256SUMS)
