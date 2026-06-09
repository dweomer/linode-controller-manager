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

## Remove the redundant Scheme from all reconciler structs.

