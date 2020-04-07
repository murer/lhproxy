#!/bin/bash -xe

cmd_cleanup() {
  docker rm -f a1 a2 || true
  docker network rm lhproxy-it-http || true
  true
}

cmd_run() {
  cmd_cleanup
  docker network create --driver bridge lhproxy-it-http --label lhproxy_dev
  docker run -dit --name a1 --network lhproxy-it-http alpine:3.11 sh
  docker run -dit --name a2 --network lhproxy-it-http alpine:3.11 sh
  docker exec -it a1 ping a2
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
