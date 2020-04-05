#!/bin/bash -xe

cmd_test() {
  go test ./pipe ./server ./util ./util/queue ./test ./cmd "$@"
}

cmd_fmt() {
  # set +x
  find -name "*.go" | grep -v "\.git" | \
    while read k; do dirname "$k"; done | sort | uniq | \
    while read k; do go fmt -x "$k" ; done
  # set -x
}


cmd_fmt2() {
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
