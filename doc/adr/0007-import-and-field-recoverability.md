# ADR-0007: Import Semantics and Field Recoverability

## Status

Accepted

## Context

The K8s CRD type does not necessarily mirror the Linode API type or operation constraints, even though one should
ideally be coercible into the other.
A key use case is **import**: a user creates a CR with only an `id` in the spec, signaling the controller to perform
a GET against the Linode API and hydrate the CR from the backend state.

This means the controller must be able to reconstruct nearly the entire CR from a read operation.
However, some fields are **write-only** in the Linode API — they were provided at creation time but are never returned
in read responses (e.g., `rootPassword`, `authorizedKeys`, `stackscriptData`).
These values are unrecoverable after creation.

## Decision

- An Instance CR with only an `id` in the spec implies an import: the controller fetches the backend state and
  populates the CR accordingly.
- Fields that can be read from the backend are populated into the CR's spec and status as appropriate.
- Fields that **cannot** be recovered from the backend (write-only fields) are noted on the runtime CR via a label
  and/or annotation indicating which values are unrecoverable.
  This makes it visible to the user that the imported CR is missing certain creation-time inputs.
- The K8s type is the controller's working model of the resource.
  It must accommodate both "created via CR" and "imported from existing backend resource" workflows.
  API operation constraints (which fields are required for create vs. which are returned on read) inform the
  translation layer but do not dictate the CRD schema.

## Consequences

- The `api:"writeonly"` (or equivalent) tag option becomes meaningful at runtime, not just for codegen — the
  controller uses it to identify which fields to flag as unrecoverable after import.
- Import is a first-class operation, not an afterthought.
  The reconciler must handle the case where spec fields beyond `id` are empty and treat it as an import signal.
- Users can adopt existing Linode infrastructure into Kubernetes management without recreating resources.
- Imported CRs are honest about what they don't know — unrecoverable values are surfaced, not silently omitted.
