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

cmd_tty() {
  LHPROXY_DOCKER_EXTRA=-it cmd_run "$@"
}

cmd_test() {
  LHPROXY_DOCKER_EXTRA=-i cmd_run go test \
    ./pipe ./server ./util ./util/queue ./test ./cmd "$@"
}

cmd_curl_test() {
  curl -v \
     --no-buffer \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
     --header "Sec-WebSocket-Protocol: chisel-v3" \
     --header "Sec-WebSocket-Version: 13" \
     http://t1.test.serasa.a.vpn.dextra.com.br/
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
