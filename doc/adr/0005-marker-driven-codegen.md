# ADR-0005: Marker-Driven Code Generation for API Translation

## Status

Proposed

## Context

ADR-0004 establishes that `api` struct tags on fields indicate which Linode API verb-face each field participates in.
The controller needs translation functions to convert between K8s types and Linode API request/response bodies — e.g.,
building a create request from only the fields without `api:"readonly"`, or populating status from only the fields
with `api:"status"`.

Writing these translation functions by hand is tedious, error-prone, and drifts as fields are added.
The `api` tags already encode the information needed to generate them.

## Decision

Deferred.
Design a code generator that consumes `api` struct tags (and the `linode` struct tag for field name mapping) to
produce typed translation functions:

- `toCreateRequest()` — includes fields without `readonly`
- `toUpdateRequest()` — includes mutable fields only (neither `readonly` nor immutable-after-create)
- `fromReadResponse()` — populates all fields from a GET response
- `statusFromReadResponse()` — populates only `status` fields into `.status`

## Consequences

- To be determined once implementation is attempted.
- Open questions: should this be a controller-gen plugin, a standalone generator, or a `go generate` tool?
  How to express "immutable after creation" distinctly from "readonly always"?
