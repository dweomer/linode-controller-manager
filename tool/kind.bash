#!/usr/bin/env bash
return 1 2>/dev/null # guard against sourcing
exec go tool "-modfile=tool/kind/go.mod" "kind" "$@"
