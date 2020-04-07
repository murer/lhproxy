#!/bin/bash -xe

cmd_cleanup() {
  docker rm -f lhproxy-it-nginx lhproxy-it || true
  docker network rm lhproxy-it-nginx || true
  true
}

cmd_run() {
  cmd_cleanup
  docker network create --driver bridge lhproxy-it-nginx --label lhproxy_dev
  docker run -dit --rm --network lhproxy-it-nginx --name lhproxy-it-nginx nginx:1.17
  docker run -it --rm --network lhproxy-it-nginx --name lhproxy-it curlimages/curl:7.69.1 curl --f http://lhproxy-it-nginx/
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
