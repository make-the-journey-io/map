#!/bin/bash -e
target_prefix="$(git rev-parse --show-toplevel)/docs/map"
go run .
go run . -graph > "$target_prefix.dot"
dot "$target_prefix.dot" -Tsvg -o"$target_prefix.svg"
git status --short "$target_prefix*"
