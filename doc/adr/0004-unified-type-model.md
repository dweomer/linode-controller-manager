# ADR-0004: Unified Type Model and Translation Layer

## Status

Accepted

## Context

The Kubernetes API expects fixed nouns — an `Instance` type is the same struct for all operations (get, create,
update, watch).
The Linode API presents different shapes of the same resource depending on the verb:

- **Create** (POST) accepts creation-time fields: region, type, image, root_pass, authorized_keys, etc.
- **Read** (GET) returns observed state including readonly/computed fields: id, status, created, ipv4, specs, alerts,
  etc.
- **Update** (PUT) accepts only a mutable subset: label, tags, alerts, backups_enabled, etc.
  Fields like region and type are immutable post-creation.

Critically, the GET response is NOT the superset of all fields.
Many fields are **write-only** — they exist only in create or update requests and never appear in read responses:

- Security-sensitive: `root_pass`
- Provisioning inputs: `authorized_keys`, `authorized_users`, `stackscript_id`, `stackscript_data`, `backup_id`,
  `swap_size`, `booted`, `metadata.user_data`
- Side-effect triggers: `private_ip`, `firewall_id`, `network_helper`
- Represented differently in read: `backups_enabled` (GET returns `backups.enabled` in a nested object instead)

This pattern holds across Linode resource types — no single API verb returns the complete picture.

The `linodego` SDK models this explicitly with separate structs per verb: `Instance` (GET response),
`InstanceCreateOptions` (POST body), `InstanceUpdateOptions` (PUT body).
A Kubernetes CRD cannot do this — it must represent all fields across all operations as a single type.
The controller acts as a translation layer between the Kubernetes API (fixed noun, standard verbs) and the Linode API
(varying noun shape per verb).

## Decision

Each Kubernetes resource type captures the full union of fields across all Linode API verb-faces for that resource.
Two struct tags work together to describe the mapping:

- **`linode:"field_name"`** — maps the K8s field to the Linode API JSON field name.
  Stays clean as a simple pointer to the backend field, readable by both humans and software.
- **`api:"opt1,opt2,..."`** — carries Linode API semantics as a comma-separated list of properties:
  - `required` — field must be provided on create.
  - `readonly` — field is set by Linode, never writable by the user (e.g., `id`, `created`).
  - `immutable` — field is writable on create but locked after (e.g., `region`, `type`, `image`).
  - `writeonly` — field is accepted on create/update but never returned in read responses (e.g., `root_password`,
    `authorized_keys`). These are unrecoverable on import (see ADR-0007).
  - `status` — field is observed state, populated from the Linode API read response.
  - `filterable` — field can be used in Linode API list filtering.
  - `deprecated` — field is deprecated in the Linode API.

Fields without an `api` tag have no special Linode API semantics (plain writable create/update fields).

Example from `api/v1alpha1/instance_type.go`:
```go
ID     int64  `json:"id,omitzero" linode:"id" api:"required,filterable,readonly,status"`
Label  string `json:"label,omitzero" linode:"label" api:"required,filterable,status"`
Region string `json:"region,omitzero" linode:"region" api:"required,filterable,readonly,status"`
Booted bool   `json:"booted" linode:"booted" default:"true"`
```

## Consequences

- The `linode` tag is a clean, single-purpose field-name pointer — never cluttered with behavioral metadata.
- The `api` tag is machine-parseable via `reflect.StructTag.Get("api")` for code generation and runtime translation
  logic.
- The `linodego` SDK is a convenience, not a dependency on the translation design.
  If linodego's types diverge from what the translation layer needs, it can be forked or bypassed — the struct tags
  are the source of truth for field mapping.
- Spec contains both user-intent fields and identity/immutable fields (like `id`, `region`) because the K8s type
  represents "the thing we are managing" in full, not just "what the user can change."
  The `api` tag disambiguates ownership.
- Future tooling (code generators, linters) can consume the `api` tag to generate translation code, validate field
  usage in reconcilers, or produce documentation.
