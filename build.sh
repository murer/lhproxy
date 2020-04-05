#!/bin/bash -xe

find_dirs_by_file() {
  set +x
  find "${1?'base dir for find, use . for all'}" -name "${2?'pattern, like: *.go'}" | grep -v "\.git" | \
    while read k; do dirname "$k"; done | sort | uniq
  set -x
}

cmd_test() {
  findbase="${1?"path is required, may be ."}"
  shift
  find_dirs_by_file "$findbase" '*_test.go' | xargs go test "$@"
  #go test ./pipe ./server ./util ./util/queue ./test ./cmd "$@"
}

cmd_fmt() {
  find_dirs_by_file '*.go' | while read k; do go fmt -x "$k" ; done
}

cmd_sshtest() {
  ssh localhost whoami
  ssh -o "ProxyCommand ./build.sh runi go run main.go client pipe native %h:%p" localhost whoami
  ssh -o "ProxyCommand ./build.sh runi go run main.go client pipe lhproxy http://localhost:8080/ %h:%p" localhost whoami
  echo SUCCESS
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
