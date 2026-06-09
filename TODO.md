# TODO

## Instance CRD — OpenAPI Spec Audit (2026-06-07)

Verified Instance CRD against cached OpenAPI spec (`.cache/openapi.json`). Locks field confirmed correct. Open items:

- [ ] **`site_type`** — Present on GET response in spec. We removed it ("belongs on Region"). May need to add back to
  InstanceStatus as readonly. Verify whether it carries meaningful per-instance info or is just inherited from region.

- [ ] **`disk_encryption`** — Spec shows GET/POST/PUT. We have `api:"immutable"`. Confirm PUT doesn't actually allow
  changing it (likely response echo, not writable).

- [ ] **`image`** — Same as disk_encryption: GET/POST/PUT in spec, `api:"immutable"` in CRD. Confirm PUT behavior.

- [ ] **`capabilities`**, **`has_user_data`** — On GET and PUT in spec. We have these as readonly in status. Likely PUT
  response echo. Low priority.

> The PUT body in the OpenAPI spec may share the response schema, making fields _appear_ writable when they're actually
> read-only in the response. Test with actual PUT requests if uncertain.

## Token Management / Auth Model

Current: single admin PAT in a Secret in the controller namespace, shared across all reconcilers and the event poller.

- [ ] **Read-only event poller token** — At startup, if the admin token has sufficient privileges, create a scoped
  read-only token for the event poller instead of sharing the admin PAT.
- [ ] **Per-user token delegation** — Controller could create limited-duration tokens for child users of the Linode
  account, rotating as necessary.
- [ ] **One connection per namespace** — The controller is namespace-scoped; each watched namespace should have its own
  Linode API connection sourced from a Secret in that namespace.
