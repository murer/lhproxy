#!/bin/bash -e

find_dirs_by_file() {
  find "${1?'base dir for find, use . for all'}" -name "${2?'pattern, like: *.go'}" | \
    grep -v "\.git" | while read k; do dirname "$k"; done | sort | uniq
}

cmd_clean() {
  rm -rf build || true
}

cmd_test() {
  findbase="${1?"path is required, may be ."}"
  shift
  find_dirs_by_file "$findbase" '*_test.go' | xargs go test "$@"
}

cmd_vendor() {
  go mod vendor -v
}

cmd_fmt() {
  find_dirs_by_file "." '*.go' | while read k; do go fmt "$k" ; done
}

cmd_build() {
  lhproxy_goos="${1?'use: linux, darwin or windows'}"
  lhproxy_goarch="${2:-"amd64"}"
  lhproxy_version="${3:-"dev"}"
  lhproxy_excname="lhproxy"
  if [[ "x$lhproxy_goos" == "xwindows" ]]; then lhproxy_excname="lhproxy.exe"; fi
  lhproxy_ldflags="-s -w -extldflags '-static' -X main.Version=$lhproxy_version"
  rm -rf "build/out/$lhproxy_goos-$lhproxy_goarch" || true
  rm "build/pack/lhproxy-$lhproxy_goos-$lhproxy_goarch-$lhproxy_version.tar.gz" || true
  mkdir -p build/pack
  CGO_ENABLED="0" GOOS="$lhproxy_goos" GOARCH="$lhproxy_goarch" \
    go build -a -trimpath -ldflags "$lhproxy_ldflags" \
      -installsuffix cgo -tags netgo -mod mod \
      -o "build/out/$lhproxy_goos-$lhproxy_goarch/lhproxy/$lhproxy_excname" .
  du -hs "build/out/$lhproxy_goos-$lhproxy_goarch/lhproxy/$lhproxy_excname"
  cd "build/out/$lhproxy_goos-$lhproxy_goarch"
  tar cvzf "../../pack/lhproxy-$lhproxy_goos-$lhproxy_goarch-$lhproxy_version.tar.gz" "lhproxy"
  cd -
}

cmd_build_all() {
  lhproxy_version="${1?"version"}"
  rm -rf "build/out" "build/pack" || true
  cmd_build windows amd64 "$lhproxy_version"
  cmd_build darwin amd64 "$lhproxy_version"
  cmd_build linux amd64 "$lhproxy_version"
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
