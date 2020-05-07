#!/bin/bash -e
go run .
go run . -graph | dot -Gdpi=200 -Tpng -o"$(git rev-parse --show-toplevel)/docs/map.png"
