#!/bin/bash -xe

cmd_detect_version() {
  echo "Branch: $TRAVIS_BRANCH"
  echo "Tag: $TRAVIS_TAG"
}

cmd_build() {
  ./docker.sh runi golang ./build.sh test .
  ./docker.sh runi golang ./build.sh build_all dev
  [[ -z "$(git status --porcelain)" ]]
}

cmd_fmt() {
  ./docker.sh runi golang ./build.sh fmt
  [[ -z "$(git status --porcelain)" ]]
}

cd "$(dirname "$0")/.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
