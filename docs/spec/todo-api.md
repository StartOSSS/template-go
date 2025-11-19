# Todo API Specification

Version: v2024-05-06 (AIP-3)

The Todo API follows spec-driven development principles: update this document and the
`proto/todo/v1/todo.proto` contract before touching implementation details.

## Overview

- Protocol: [Connect](https://connect.build) over HTTP/2 or HTTP/1.1
- Payloads: Protobuf + JSON interoperability
- Authentication: pluggable via Cloud Run IAM or Envoy filters (not part of the sample)

### Planes (AIP-111)

- **Management plane:** Helm charts, Terraform modules, and Skaffold orchestrate deployments, policy, and observability.
- **Data plane:** The `TodoService` RPCs, Connect handlers, and PostgreSQL store operate on user todo items.

## RPCs

| RPC | Description |
|-----|-------------|
| `CreateTodo` | Create a new todo task |
| `GetTodo` | Fetch a task by UUID |
| `ListTodos` | Enumerate tasks ordered by recency |
| `UpdateTodo` | Update metadata and completion flag |
| `DeleteTodo` | Remove a task |
| `HealthCheck` | Basic readiness probe |

Refer to the proto file for field-level semantics and the [AIP glossary](https://google.aip.dev/9) for terminology.

## Versioning & review checkpoints

- The proto namespace `todo.v1` follows [AIP-2](https://google.aip.dev/2) numbering conventions.
- This spec's `Version:` header follows [AIP-3](https://google.aip.dev/3); bump the date whenever an RPC, field, or behavior change is proposed.
- Before beta/GA releases, copy the API design review checklist from [AIP-100](https://google.aip.dev/100) into the PR description and document reviewer approval links here.

## Persistence

- PostgreSQL 15+
- Table: `todo_tasks`
- Columns: `id uuid`, `title text`, `description text`, `completed bool`, `created_at timestamptz`, `updated_at timestamptz`

## Non-functional requirements

- P99 latency < 200ms at 50 RPS in local load tests
- Emit OTel spans + metrics for every RPC
- Health checks must fail if the DB connection is unhealthy
- Tests: unit (`go test`), integration (black-box client), load (k6)

## Beta blockers (AIP-205)

- [ ] Document final IAM story for multi-tenant deployments (tracking issue TBD)
- [ ] Validate pagination ergonomics for `ListTodos`

## Precedent tracking (AIP-200)

Add inline `TODO(aip-200)` annotations inside the Go code and mirror them in the PR description if a temporary exception to any guidance is necessary.

## Google AIP alignment resources

See `docs/spec/aip-compliance.md` for the living checklist that applies to this spec and every derivative service.
