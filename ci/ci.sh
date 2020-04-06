#!/bin/bash -xe

cmd_detect_version() {
  mkdir build
  LHPROXY_VERSION="$TRAVIS_TAG"
  if [[ -z "$LHPROXY_VERSION" ]]; then
    LHPROXY_VERSION="branch-$TRAVIS_BRANCH"
  fi
  echo "$LHPROXY_VERSION" > build/version.txt
}

cmd_build() {
  ./docker.sh runi golang go version
  ./it/it.sh build base 1> /dev/null 2>&1 &
  ./docker.sh runi golang ./build.sh test .
  ./docker.sh runi golang ./build.sh build_all "$LHPROXY_VERSION"
   wait %1
  ./it/it.sh it
  ./docker.sh build
  [[ -z "$(git status --porcelain)" ]]
}

cmd_fmt() {
  ./docker.sh runi golang ./build.sh fmt
  [[ -z "$(git status --porcelain)" ]]
}

cmd_deploy_docker() {
  set +x
  docker login --username "$DOCKERHUB_USER" --password "$DOCKERHUB_PASS"
  set -x
  LHPROXY_DOCKER_VERSION="$(echo "$LHPROXY_VERSION" | cut -d'-' -f2-)"
  ./docker.sh push "$LHPROXY_DOCKER_VERSION"
  ./docker.sh push latest
}

cd "$(dirname "$0")/.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
