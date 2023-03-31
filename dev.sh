#!/usr/bin/env bash

set -e

figwheel() {
  clojure -M -m figwheel.main "$@"
}

build() {(
  figwheel -bo dev
)}

build_release() {
  echo "Not implemented yet."
  exit 1
}

run() {
  figwheel -b dev -r
}

backend() {(
  cd backend; go run ./cmd/server
)}

show_build_config() {
  figwheel -pc -bo dev
}

cmd="${1:-build}"
shift || true
$cmd "$@"
