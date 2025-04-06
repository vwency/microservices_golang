#!/bin/bash

if [ -f go.mod ]; then
  awk '/^require \(/,/^\)/' go.mod | \
  grep -v '^require (' | \
  grep -v '^)' > go.mod.deps
  echo "Dependencies updated in go.mod.deps"
else
  echo "Error: go.mod not found!"
  exit 1
fi
