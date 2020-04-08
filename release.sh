#!/bin/bash -xe

cmd_version() {
  lhproxy_version="${1?'version, like: x.x.x'}"
  echo "$lhproxy_version" | grep "^[0-9]\+\.[0-9]\+\.[0-9]\+$"
  git tag "$lhproxy_version"
  git push origin "$lhproxy_version"
}

cmd_edge() {
  cmd_force edge
}

cmd_force() {
  lhproxy_version="${1?'version, like: x.x.x'}"
  git tag "$lhproxy_version" -f
  git push origin "$lhproxy_version" -f
}

cmd_DeLeTe_TaG() {
  lhproxy_version="${1?'version, like: x.x.x'}"
  git tag -d "$lhproxy_version"
  git push --delete origin "$lhproxy_version"
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
