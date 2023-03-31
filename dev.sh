#!/usr/bin/env bash

set -e

build() {
  (cd backend; go build -o ~/bin/noport ./cmd/noport)
  # figwheel -bo dev
}

fmt() {
  (cd backend; go fmt ./...)
}

figwheel() {
  clojure -M -m figwheel.main "$@"
}

front() {
  figwheel -b dev -r
}

figwheel_config() {
  figwheel -pc -bo dev
}

cmd="${1:-build}"
shift || true
$cmd "$@"
