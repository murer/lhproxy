#!/bin/bash -xe

cmd_cleanup() {
  docker rm -f lhproxy_it_squid || true
  docker rm -f lhproxy_it_pipe || true
}

cmd_build() {
  lhproxy_it_target="${1:-final}"
  docker build -t lhproxy/it:dev --target "$lhproxy_it_target" .
}

cmd_run() {
  trap cmd_cleanup EXIT
  cmd_cleanup
  cmd_build
  docker run -d --rm --label lhproxy_dev --name lhproxy_it_squid \
    -p 3128:3128 -h lhproxy_it_squid lhproxy/it:dev /root/entrypoint/server.sh

  docker run -it --rm --label lhproxy_dev --name lhproxy_it_pipe \
    --network host lhproxy/it:dev /root/entrypoint/pipe.sh

  echo "SUCCESS"
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
