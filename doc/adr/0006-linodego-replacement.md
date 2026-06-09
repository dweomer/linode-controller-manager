# ADR-0006: Struct Tags as Foundation for a linodego Replacement

## Status

Deferred

## Context

The `linode` and `api` struct tags (ADR-0004) capture field-level translation metadata between K8s types and Linode
API request/response bodies.
The question is whether this tag system can serve as the foundation for replacing linodego entirely — generating or
building a Linode API client directly from the struct metadata.

## What the tags already cover

- Field name mapping to Linode API JSON fields (`linode` tag)
- Which fields appear in read responses (`api:"status"`)
- Which fields are user-writable vs Linode-assigned (`api:"readonly"`)
- Which fields are required on create (`api:"required"`)
- Which fields support list filtering (`api:"filterable"`)

## Known gaps

### Field-level

- **Immutable vs readonly** — `api:"readonly"` conflates "never set by user" (`id`, `created`) with "set on create
  but locked after" (`region`, `type`). Need something like `immutable` to distinguish these.
- **Write-only** — fields like `rootPassword`, `authorizedKeys` are submitted on create but never returned.
  Currently implicit (no `status` or `readonly`), but an explicit `writeonly` would be clearer.
- **Shape mismatches** — `backups_enabled` on create maps to `backups.enabled` in the read response.
  The `linode` tag points to one path; there's no way to express differing paths per verb.

### Resource-level

- API endpoint paths (`/v4/linode/instances/{id}`)
- Supported verbs per resource (not all resources support full CRUD)
- Action endpoints (boot, reboot, shutdown, resize, clone, migrate, rebuild)
- Sub-resource relationships (instance → disks, configs, IPs)
- Pagination strategy for list operations

### Client-level

- Auth, base URL, retries, rate limiting, error model

## Decision

Deferred until more Linode API types are implemented.
The Instance type alone may not surface all the patterns needed to design a complete replacement.
Revisit after implementing at least a few more resource types to see which gaps are real and which are
Instance-specific.
