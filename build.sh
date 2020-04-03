#!/bin/bash -xe

cmd_golang() {
  docker volume create lhproxy_golang_dev --label lhproxy_dev || true
  docker run -it --rm --label lhproxy_dev \
    --mount source=lhproxy_golang_dev,target=/go \
    -v "$(pwd)":/go/src -w /go/src \
    -p 8080:8080 \
    golang:1.14 "$@"
}

cmd_test() {
  # --proxy http://C80795A:QWEasd1@@localhost:6001
  cmd_golang go test ./util -v
}

cmd_curl_test() {
  #--proxy-pass 'QWEasd1@' --proxy-user C80795A --proxy-basic --proxy http://localhost:6001/ \
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

# GET http://localhost:6001/ HTTP/1.1
# Host: t1.test.serasa.a.vpn.dextra.com.br
# User-Agent: Go-http-client/1.1
# Connection: Upgrade
# Proxy-Authorization: Basic QzgwNzk1QTpRV0Vhc2QxQA==
# Sec-WebSocket-Key: FYtltTOZX/HBH2iT1BQxgQ==
# Sec-WebSocket-Protocol: chisel-v3
# Sec-WebSocket-Version: 13
# Upgrade: websocket
