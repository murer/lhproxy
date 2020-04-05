#!/bin/bash -xe

LHPROXY_USER_ID="$(id -u):$(id -g)"

docker_hub() {
  docker run $LHPROXY_DOCKER_EXTRA --rm --label lhproxy_dev \
    lhproxy/hub:local "$@"
}

docker_golang() {
  docker volume create lhproxy_golang_dev --label lhproxy_dev || true
  docker run $LHPROXY_DOCKER_EXTRA --rm --label lhproxy_dev \
    --mount source=lhproxy_golang_dev,target=/go \
    -v "$(pwd)":/go/src -w /go/src \
    --network host \
    -e "LHPROXY_SECRET=12345678901234561234567890123456" \
    -u "$LHPROXY_USER_ID" \
    -e "HOME=/go" \
    golang:1.14 "$@"
}

cmd_images() {
  cd docker
  docker build -t lhproxy/hub:local -f Dockerfile.hub .
  cd -
}

cmd_clean() {
  ./build.sh clean
  docker ps -aq --filter label lhproxy_dev | xargs docker rm -f || true
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
