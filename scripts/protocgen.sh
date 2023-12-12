#!/usr/bin/env bash

# see pre-requests:
# - https://grpc.io/docs/languages/go/quickstart/
# - gocosmos plugin is automatically installed during scaffolding.

set -eo pipefail
go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null

echo "Generating gogo proto code"
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package "$file" ; then
      buf generate --template buf.gen.gogo.yaml "$file"
    fi
  done
done


# move proto files to the right places
cp -r github.com/CudoVentures/cudos-node/* ./
rm -rf github.com

