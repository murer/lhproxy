#!/bin/bash -xe

./build.sh test

[[ -z "$(git status --porcelain)" ]]

./build.sh fmt

[[ -z "$(git status --porcelain)" ]]
