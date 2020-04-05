#!/bin/bash -xe

find_dirs_by_file() {
  set +x
  find "${1?'base dir for find, use . for all'}" -name "${2?'pattern, like: *.go'}" | \
    grep -v "\.git" | while read k; do dirname "$k"; done | sort | uniq
  set -x
}

cmd_test() {
  findbase="${1?"path is required, may be ."}"
  shift
  find_dirs_by_file "$findbase" '*_test.go' | xargs go test "$@"
}

cmd_fmt() {
  find_dirs_by_file "." '*.go' | while read k; do go fmt -x "$k" ; done
}

cmd_build() {
  lhproxy_goos="${1?'use: linux, darwin or windows'}"
  lhproxy_goarch="${2:-"amd64"}"
  lhproxy_version="${3:-"dev"}"
  lhproxy_excname="lhproxy"
  if [[ "x$lhproxy_goos" == "xwindows" ]]; then lhproxy_excname="lhproxy.exe"; fi
  lhproxy_ldflags="-s -w -extldflags '-static'"
  rm -rvf "build/out/$lhproxy_goos-$lhproxy_goarch" || true
  CGO_ENABLED="0" GOOS="$lhproxy_goos" GOARCH="$lhproxy_goarch" \
    go build -a -trimpath -ldflags "$ldflags" \
      -installsuffix cgo -tags netgo -mod mod \
      -o "build/out/$lhproxy_goos-$lhproxy_goarch/lhproxy-$lhproxy_version/$lhproxy_excname" .
}

cmd_build_all() {
  lhproxy_version="${1?"version"}"
  rm -rf "build/out" || true
  cmd_build windows amd64
  cmd_build darwin amd64
  cmd_build linux amd64
}

cmd_sshtest() {
  ssh localhost whoami
  ssh -o "ProxyCommand ./build.sh runi go run main.go client pipe native %h:%p" localhost whoami
  ssh -o "ProxyCommand ./build.sh runi go run main.go client pipe lhproxy http://localhost:8080/ %h:%p" localhost whoami
  echo SUCCESS
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
