#!/usr/bin/env bash

set -e

bin_default=~/bin/noport

build() {
  build_front
  build_back "$1"
}

build_back() {
  # Assumes we just did build_front.
  rm -rf backend/public
  cp -r resources/public backend/public
  cp -r target/public/js backend/public/js
  (cd backend; go build -o "${1:-$bin_default}" ./cmd/noport)
}

build_front() {
  figwheel -bo prod
}

fmt() {
  (cd backend; go fmt ./...)
}

figwheel() {
  export NOPORT_DEV=1
  clojure -M -m figwheel.main "$@"
}

server() {
  export NOPORT_DEV=1
  build_back "$1"
  noport
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
