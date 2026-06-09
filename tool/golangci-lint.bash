#!/usr/bin/env bash
return 1 2>/dev/null # guard against sourcing
exec go tool "-modfile=tool/golangci-lint-kube-api-linter/go.mod" "golangci-lint-kube-api-linter" "$@"
