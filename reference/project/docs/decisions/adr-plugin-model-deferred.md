---
namespace: feather
type: adr
id: plugin-model-deferred
instance-version: 1
revision: 1
content-state: accepted
existence-state: active
dimension: decisions
author: Lukas Weber
created: 2026-06-20
updated: 2026-06-26
supersedes: []
derives-from: []
depends-on:
  - req:publishing-core
amends: []
validates: []
change-log:
  - date: 2026-06-20
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-06-20
    domain: content-state
    from: "-"
    to: proposed
    by: Lukas Weber
  - date: 2026-06-26
    domain: content-state
    from: proposed
    to: accepted
    by: Lukas Weber
---

# ADR — No Plugin Model in v1

## Context

The vision promises "no plugin runtime" and the 2026 strategy explicitly defers extensibility. Still, during architecture design the question recurs: should the rendering or API layer reserve extension points for themes/plugins, since retrofitting one later is expensive?

## Decision

Ship **no plugin model in v1**. No extension points in the renderer or the API beyond plain configuration (theme directory, site settings). The strategy's "extensibility is a future problem" is accepted as a constraint on the architecture.

## Consequences

- The renderer stays a pure function over (template, markdown); there is no plugin loading, sandboxing, or versioning surface to secure.
- Engineering budget is preserved for the authoring loop, which the strategy ranks first.
- Known cost: if plugins are ever added, the render pipeline may need reshaping; the risk is accepted and bounded because the pipeline is small and internal-only.
- Configuration is YAML, not code: customization in v1 means editing templates, which keeps the single-binary property.

## Alternatives Considered

- **Pluggable renderer now** — rejected: a plugin ABI, sandbox, and lifecycle are a second product; no demand signal exists.
- **Theme-only extension** — accepted in part: templates are swappable via configuration, which covers most "plugin" requests without a runtime.
