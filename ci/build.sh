#!/bin/bash -xe

cmd_script() {
  ./build.sh test
  [[ -z "$(git status --porcelain)" ]]
  ./build.sh fmt
  [[ -z "$(git status --porcelain)" ]]
}

cd "$(dirname "$0")/.."; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
