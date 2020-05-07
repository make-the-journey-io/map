#!/bin/bash -e
target="$(git rev-parse --show-toplevel)/docs/map.png"
go run .
go run . -graph | dot -Gdpi=200 -Tpng -o"$target"
git status --short "$target"
