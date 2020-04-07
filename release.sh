#!/bin/bash -xe

cmd_version() {
  lhproxy_version="${1?'version, like: x.x.x'}"
  echo "$lhproxy_version" | grep "^[0-9]\+\.[0-9]\+\.[0-9]\+$"
  echo git tag "$lhproxy_version"
  echo git push origin "$lhproxy_version"
}

cmd_edge() {
  git tag edge
  git push origin edge
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
