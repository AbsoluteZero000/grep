#!/bin/sh

set -e # Exit on failure

go build -o /tmp/grep-go cmd/mygrep/main.go
