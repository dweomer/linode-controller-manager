# ADR-0002: Namespace-Scoped Resources

## Status

Accepted

## Context

Kubernetes CRDs can be either namespace-scoped or cluster-scoped.
This choice affects RBAC granularity, multi-tenancy, and operational complexity.

Cluster-scoped resources (like Crossplane's model) allow a single controller instance to manage all infrastructure,
but make multi-tenant isolation harder — RBAC can only grant access to the entire resource type, not subsets.

Namespace-scoped resources allow natural tenant isolation via namespace RBAC, simpler least-privilege policies, and
familiar Kubernetes workflows (namespace-per-team, namespace-per-environment).

## Decision

All Linode resources are namespace-scoped.
The controller watches only namespaces specified in the `WATCH_NAMESPACE` environment variable (required,
comma-separated). RBAC is generated as a `Role`, not `ClusterRole`.

## Consequences

- Multi-tenancy works via standard namespace isolation — teams get RBAC on their namespace, can't see or modify other
  teams' Linode resources.
- The controller can be deployed per-namespace or watching multiple namespaces, depending on the operator's needs.
- A single Linode resource cannot span namespaces (acceptable — Linode resources are independent entities).
- `WATCH_NAMESPACE` is required at startup; the controller will not default to watching all namespaces.
