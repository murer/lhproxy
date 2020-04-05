#!/bin/bash -xe

cmd_detect_version() {
  export LHPROXY_VERSION="$TRAVIS_TAG"
  if [[ -z "$LHPROXY_VERSION" ]]; then
    export LHPROXY_VERSION="branch-$TRAVIS_BRANCH"
  fi
  export LHPROXY_VERSION
  export GITHUB_TAG_NAME="$LHPROXY_VERSION"
}

cmd_build() {
  # LHPROXY_VERSION="$(cat build/version.txt)"
  ./docker.sh runi golang ./build.sh test .
  ./docker.sh runi golang ./build.sh build linux amd64 "$LHPROXY_VERSION"
  [[ -z "$(git status --porcelain)" ]]
}

cmd_fmt() {
  ./docker.sh runi golang ./build.sh fmt
  [[ -z "$(git status --porcelain)" ]]
}

cd "$(dirname "$0")/.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
