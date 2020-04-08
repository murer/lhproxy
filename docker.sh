#!/bin/bash -e

LHPROXY_USER_ID="$(id -u):$(id -g)"

cmd_init() {
  docker network create --label lhproxy_dev --driver bridge lhproxy-dev-network
}

cmd_build() {
  docker build --target lhproxy_scratch -t lhproxy/lhproxy:dev .
  docker build --target lhproxy_alpine -t lhproxy/lhproxy:dev-alpine .
}

cmd_push() {
  lhproxy_docker_version="${1?"version to push"}"
  docker tag lhproxy/lhproxy:dev "murer/lhproxy:$lhproxy_docker_version"
  docker tag lhproxy/lhproxy:dev-alpine "murer/lhproxy:$lhproxy_docker_version-alpine"
  docker push "murer/lhproxy:$lhproxy_docker_version"
  docker push "murer/lhproxy:$lhproxy_docker_version-alpine"
}

docker_lhproxy() {
  docker run $LHPROXY_DOCKER_EXTRA --rm --label lhproxy_dev \
    -u "$LHPROXY_USER_ID" lhproxy/lhproxy:dev "$@"
}

docker_golang() {
  docker volume create lhproxy_golang_dev --label lhproxy_dev || true
  docker run $LHPROXY_DOCKER_EXTRA --rm --label lhproxy_dev \
    --network lhproxy-dev-network \
    --mount source=lhproxy_golang_dev,target=/go \
    -v "$(pwd)":/go/src -w /go/src \
    -e "LHPROXY_SECRET=123" \
    -e "HOME=/go" \
    -u "$LHPROXY_USER_ID" \
    golang:1.14 "$@"
}

cmd_cleanup() {
  ./build.sh clean
  docker ps -aq --filter label=lhproxy_dev | xargs docker rm -f || true
  docker system prune --volumes --filter label=lhproxy_dev -f || true
}

cmd_run() {
  dockername="${1?'docker name'}"
  shift
  "docker_${dockername}" "$@"
}

cmd_runi() {
  istty=-i
  [[ -t 0 ]] && istty=-it
  LHPROXY_DOCKER_EXTRA="$LHPROXY_DOCKER_EXTRA $istty" cmd_run "$@"
}

cmd_rund() {
  LHPROXY_DOCKER_EXTRA="$LHPROXY_DOCKER_EXTRA -d" cmd_run "$@"
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
