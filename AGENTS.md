# Agents

## controller

Implements and evolves reconcilers for Linode resources following the ACK pattern: spec is user intent, status
reflects observed Linode API state.

### Scope

- Reconcile loop: fetch CR, call Linode API via linodego, update status
- Finalizer management for resource deletion
- Status condition updates using `metav1.Condition`
- Error handling, requeue strategies, backoff on rate limits
- Event recording for user-visible actions

### Key Files

- `internal/controller/linode/instance.go`
- `api/v1alpha1/instance.go` (spec/status contract)
- `main.go` (controller registration)

### Instructions

- Use `logf.FromContext(ctx)` for structured logging
- Return `ctrl.Result{RequeueAfter: duration}` for polling Linode state
- Set a finalizer before creating external resources
- Update `.status.conditions` with standard types (Ready, Reconciling, Degraded)
- Never modify `.spec` from the controller; spec is user-owned, status is controller-owned
- Map linodego HTTP 404 to "not found" condition, 429 to requeue with backoff
- Use `client.MergeFrom` for status patch operations

## api-types

Designs and maintains CRD type definitions under `api/v1alpha1/`. Types should mirror Linode API structures.

### Scope

- Instance spec and status field design
- Kubebuilder marker annotations for validation, print columns, selectable fields
- Enum types and constants in `api/v1alpha1/instance_enums.go`
- Alignment with Linode API schema
- New resource types as the project expands beyond Instance

### Key Files

- `api/v1alpha1/instance.go`
- `api/v1alpha1/instance_types.go`
- `api/v1alpha1/instance_enums.go`
- `api/v1alpha1/group-version.go`

### Instructions

- Follow Kubernetes API conventions: `metav1.Condition` for status conditions
- Use `omitzero` json tag per project convention
- Validate fields with kubebuilder validation markers
- Keep enum types in `instance_enums.go` (or `<resource>_enums.go`) with `+k8s:enum` marker
- After any type change run `go generate ./...` and verify CRD output
- Use `+listType=set` for unordered slices, `+listType=map` with `+listMapKey` for conditions

## testing

Writes and maintains tests for controllers and API types.

### Scope

- Unit tests for reconcilers using envtest or fake client
- Integration tests with a kind cluster (`bash tool/kind.bash`)
- API type validation tests
- Test fixtures using example manifests in `manifest/example/`

### Key Files

- `internal/controller/linode/` (code under test)
- `manifest/example/instance/` (fixtures)

### Instructions

- Use `sigs.k8s.io/controller-runtime/pkg/envtest` for controller tests
- Prefer standard `testing` package over ginkgo
- Mock the Linode API client via interface; never call the real API in tests
- Test reconcile outcomes: status updates, requeue behavior, finalizer lifecycle
- Test error paths: API failures, missing resources, invalid specs
- Place test files adjacent to implementation

## manifests

Validates and evolves Kubernetes manifests and deployment configuration.

### Scope

- CRD manifest correctness (generated via `go generate`, validated separately)
- RBAC manifest review
- Example CR authoring for testing and documentation
- Deployment manifests and kustomize overlays as needed

### Key Files

- `manifest/crd/`
- `manifest/rbac/controller-manager.yaml`
- `manifest/example/instance/`

### Instructions

- Validate CRDs with `bash tool/kubeconform.bash manifest/crd/`
- RBAC is namespace-scoped; do not generate ClusterRoles
- Example CRs should exercise different field combinations and edge cases
- Deployment manifests must require `WATCH_NAMESPACE` env var in the pod spec
