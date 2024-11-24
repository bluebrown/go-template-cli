#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail
cd "$(readlink -f "$(dirname "$0")/..")"

vers="dev"
dest="dist"

while getopts "v:o:" opt; do
  case $opt in
  v) vers="$OPTARG" ;;
  o) dest="$OPTARG" ;;
  *) exit 1 ;;
  esac
done
shift $((OPTIND - 1))

if [ $# -eq 0 ] || [ "$1" = "-" ]; then
  mapfile -t targets
else
  targets=("$@")
fi

mkdir -p "$dest"

for build in "${targets[@]}"; do
  os="$(echo "$build" | cut -d'/' -f1)"
  arch="$(echo "$build" | cut -d'/' -f2)"

  echo "Building for $os/$arch"

  GOOS="$os" GOARCH="$arch" go build \
    -o "$dest/tpl-$os-$arch" \
    -ldflags "-w -s -X main.version=$vers" \
    ./cmd/tpl

done
