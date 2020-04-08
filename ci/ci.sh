#!/bin/bash -xe

cmd_detect_version() {
  mkdir -p build
  LHPROXY_VERSION="$TRAVIS_TAG"
  if [[ "x$LHPROXY_VERSION" == "x" ]]; then
    LHPROXY_VERSION="$TRAVIS_BRANCH"
  fi
  echo "$LHPROXY_VERSION" > build/version.txt
}

cmd_build() {
  ./docker.sh runi golang go version
  ./it/ssh/build.sh build base 1> /dev/null 2>&1 &
  ./docker.sh runi golang ./build.sh test .
  ./docker.sh runi golang ./build.sh build_all "$LHPROXY_VERSION"
   ./docker.sh build
   wait %1
  ./it/ssh/build.sh run

  [[ -z "$(git status --porcelain)" ]]
}

cmd_fmt() {
  ./docker.sh runi golang ./build.sh fmt
  [[ -z "$(git status --porcelain)" ]]
}

cmd_deploy_docker() {
  lhproxy_ver_alias="$1"
  set +x
  docker login --username "$DOCKERHUB_USER" --password "$DOCKERHUB_PASS"
  set -x
  ./docker.sh push "$LHPROXY_VERSION"
  if [[ "x$lhproxy_ver_alias" != "x" ]]; then
    ./docker.sh push "$lhproxy_ver_alias"
  fi
}

cd "$(dirname "$0")/.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
