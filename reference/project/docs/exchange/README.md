# docs/exchange/ — Exchange Layer (EX)

> Anchor EKA: Exchange Layer — validation + transfer. EX owns no content and no state.

## Role

The Exchange Layer manages the **boundary** of the serialization: it validates artifact conformance against the EKA contract and governs import/export between repositories. EX is the layer that "owns neither" — it owns no content (owned by the Knowledge Layer) and no state (owned by the Operating Layer).

## What EX Owns (Contract)

| Contract | Where |
|---|---|
| Conformance Rules (9 mechanical rules) | [validation.md](validation.md) |
| Import/export conventions (round-trip, identity conflicts, schema) | [transfer.md](transfer.md) |
| Token table, reference format, and identity format | serialized in the project README; EX enforces them during validation |

## What EX Does Not Own

- **Content** — artifact content belongs to the dimension owner (Knowledge Layer). EX only checks its form.
- **State** — state values and transitions belong to the state owner (Operating Layer). EX only checks value validity, never changes it.

## Usage Flow

1. A writer creates/modifies an artifact.
2. Before commit: run the [validation.md](validation.md) checklist.
3. For cross-repository exchange: follow [transfer.md](transfer.md).
4. Violations = reject commit until conformant.

## Related

- [../operating/protocol.md](../operating/protocol.md) — the rules EX validates originate in the OS protocol.
- [../README.md](../README.md) — identity, state, phase, relationship, classification, projection conventions.
