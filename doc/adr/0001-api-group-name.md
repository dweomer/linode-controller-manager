# ADR-0001: API Group Name

## Status

Accepted

## Context

This project is a Kubernetes controller manager that declares and reconciles Linode API resources as custom resources
(similar to AWS Controllers for Kubernetes). We needed to choose a Kubernetes API group name for our CRDs.

Candidates considered:

- `linode.k8s` — similar to ACK's `*.services.k8s.aws` pattern, but `.k8s` is not a real domain and could collide
  with official Kubernetes SIG projects.
- `api.linode.com` — reflects the Linode API directly, but `api.linode.com/v1alpha1` looks confusingly like a Linode
  REST API URL path (`api.linode.com/v4/instances/...`).
- `linode.com` — owned domain, short, no ambiguity with REST API URLs, Kubernetes version strings (`v1alpha1`) are
  clearly distinct from Linode API versions (`v4`, `v4beta`).

## Decision

Use `linode.com` as the API group name.

Resources are addressed as `linode.com/v1alpha1 Instance`, CRD file is `linode.com_instances.yaml`.

## Consequences

- Clear domain ownership prevents group name collisions.
- No confusion with Linode REST API URL patterns.
- If Linode ever ships an official operator using `linode.com`, we'd need to coordinate or migrate — acceptable risk
  for a community project.
