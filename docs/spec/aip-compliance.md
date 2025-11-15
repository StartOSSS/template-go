# Google AIP Alignment

Version: v2024-05-06 (AIP-3)

This template intentionally aligns with the generally applicable Google API Improvement Proposals (AIPs).
Use this checklist whenever the protobuf contract, documentation set, or runtime behavior changes.

| AIP | Topic | Template application |
|-----|-------|----------------------|
| [1](https://google.aip.dev/1) | Purpose & guidelines | All API design decisions originate from the spec-first workflow (`docs/spec/todo-api.md` + `proto/todo/v1/todo.proto`). Each change starts as documentation, mirroring the "single source of truth" directive. |
| [2](https://google.aip.dev/2) | Numbering | The proto package namespace `todo.v1` stays within the 100-999 general guidance range and new APIs must request editor-approved numbers before publishing a package version. |
| [3](https://google.aip.dev/3) | Versioning | Specs and supporting docs carry date-based tags (see the version headers in this file and `docs/spec/todo-api.md`). Release notes mention the corresponding tag whenever a significant change lands. |
| [8](https://google.aip.dev/8) | Style & guidance | Each document focuses on a discrete concern (spec, consumer, contributor, observability) and links to authoritative guidance instead of duplicating it. Protos expose descriptive comments and avoid anti-pattern-only sections. |
| [9](https://google.aip.dev/9) | Glossary | Terminology in docs matches the AIP glossary ("API consumer", "API backend", "plane", etc.). New terms must link back here so contributors share the same vocabulary. |
| [100](https://google.aip.dev/100) | API design review FAQ | The `docs/spec/todo-api.md` checklist includes review readiness questions so that every beta/GA release has a documented review thread and reviewer sign-off captured in the Git history. |
| [111](https://google.aip.dev/111) | Planes | The template explicitly calls out which features live in the management plane (Helm/Terraform/GitOps) vs. the data plane (Todo RPCs + storage). Monitoring dashboards cover both planes. |
| [200](https://google.aip.dev/200) | Precedent | If an unavoidable violation appears, open an issue and reference AIP-200 in the code comment/commit message describing the exception and why it cannot establish precedent. |
| [205](https://google.aip.dev/205) | Beta-blocking changes | The spec includes a "Beta blockers" checklist. Any open box must have a link to the tracking issue and must be resolved before tagging a beta release. |

## Working agreement

1. **Document first.** Update `docs/spec/todo-api.md` (or a sibling spec) before touching Go code or protobuf files.
2. **Tag specs with dates.** When meaningful API changes land, bump the `Version:` field in every affected spec file to the current ISO-8601 date.
3. **Track deviations.** Use `TODO(aip-200)` style comments referencing the GitHub issue when temporarily diverging from guidance.
4. **Review rigorously.** Do not promote a change set without copying the latest AIP checklist into the pull request description and confirming every box.
5. **Share vocabulary.** Link back to this file whenever adding new docs so readers inherit the glossary.

These guardrails ensure the template honors Google's recommended API practices and keeps future services consistent.
