#!/bin/bash

# Установите значения по умолчанию
GO_MODULE_PATH=${GO_MODULE_PATH:-"github.com/vwency/microservices_golang/"}
GO_VERSION=${GO_VERSION:-"1.22.0"}
GO_TOOLCHAIN=${GO_TOOLCHAIN:-"go1.24.1"}

if [ ! -f go.mod.deps ]; then
  echo "Error: go.mod.deps not found!"
  exit 1
fi

cat <<EOF > go.mod
module $GO_MODULE_PATH

go $GO_VERSION

toolchain $GO_TOOLCHAIN

require (
$(cat go.mod.deps)
)
EOF

echo "Generated go.mod with:"
echo "Module: $GO_MODULE_PATH"
echo "Go version: $GO_VERSION"
echo "Toolchain: $GO_TOOLCHAIN"
echo "Dependencies preserved from go.mod.deps"
