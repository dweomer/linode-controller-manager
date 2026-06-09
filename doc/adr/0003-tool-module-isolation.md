# ADR-0003: Tool Module Isolation

## Status

Accepted

## Context

Go 1.24+ supports declaring tool dependencies directly in `go.mod` via the `tool` directive.
This is elegant for simple projects, but Kubernetes-adjacent projects pull in massive dependency trees
(controller-runtime, client-go, apimachinery).
Dev tools like controller-gen, golangci-lint, and kubeconform bring their own large dependency graphs.

When tools share `go.mod` with the main module, their transitive dependencies participate in version resolution.
This means a linter's dependency on an older (or newer) version of a shared library can force the runtime module to
upgrade or constrain unnecessarily.
In the Kubernetes ecosystem where breaking changes across minor versions are common, letting tool needs dictate runtime
dependency versions is unacceptable.

## Decision

Each dev tool gets its own `go.mod` in `tool/<name>/go.mod` and a corresponding bash wrapper at `tool/<name>.bash`.
The main module's `go.mod` contains only runtime dependencies.

Tools are invoked via:
```bash
exec go tool "-modfile=tool/<name>/go.mod" "<binary>" "$@"
```

## Consequences

- Tool dependency upgrades never affect the runtime module's dependency resolution.
- Each tool can pin its own versions independently — controller-gen can use a different k8s library version than the
  main module without conflict.
- Slightly more files to maintain (one `go.mod` + `go.sum` per tool), but these rarely need manual attention.
- `go generate` invokes tools via the `-modfile` flag, keeping the workflow simple and reproducible.
