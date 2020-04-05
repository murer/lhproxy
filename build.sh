#!/bin/bash -xe

cmd_run() {
  docker volume create lhproxy_golang_dev --label lhproxy_dev || true
  docker run $LHPROXY_DOCKER_EXTRA --rm --label lhproxy_dev \
    --mount source=lhproxy_golang_dev,target=/go \
    -v "$(pwd)":/go/src -w /go/src \
    --network host \
    -e "LHPROXY_SECRET=12345678901234561234567890123456" \
    golang:1.14 "$@"
}

cmd_runi() {
  LHPROXY_DOCKER_EXTRA=-i cmd_run "$@"
}

cmd_runit() {
  LHPROXY_DOCKER_EXTRA=-it cmd_run "$@"
}

cmd_rund() {
  LHPROXY_DOCKER_EXTRA=-d cmd_run "$@"
}

cmd_test() {
  cmd_runi go test \
    ./pipe ./server ./util ./util/queue ./test ./cmd "$@"
}

cmd_fmt() {
  set +x
  docker rm -f lhproxy_golang_fmt || true
  find -name "*.go" | grep -v "\.git" | \
    while read k; do dirname "$k"; done | sort | uniq | \
    while read k; do \
      LHPROXY_DOCKER_EXTRA=-i cmd_run go fmt -x
    done
  set -x
}

cmd_sshtest() {
  ssh localhost whoami
  ssh -o "ProxyCommand ./build.sh runi go run main.go client pipe native %h:%p" localhost whoami
  ssh -o "ProxyCommand ./build.sh runi go run main.go client pipe lhproxy http://localhost:8080/ %h:%p" localhost whoami
  echo SUCCESS
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
