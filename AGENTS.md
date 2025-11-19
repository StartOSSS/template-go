# Repository Agent Notes

## Scope
This file applies to the entire repository.

## Guidelines
- Prefer declarative configuration (YAML/TOML) over bespoke scripting when possible.
- Go code should follow standard `gofmt`/`goimports` style and be organized using layered packages (`cmd`, `internal`, `pkg`).
- Tests should live next to the code they cover.
- Document new workflows or commands in the README when they are user-facing.
- Keep Makefile targets idempotent and safe to run multiple times.
- Default to upstream maintained building blocks (e.g. Helm charts, terraform-google-modules, OSS CLI tools) before writing custom infrastructure code.
- When adding automation, ensure it aligns with SLSA + SSDF guidance already referenced in `docs/` and surfaces clear observability hooks.
