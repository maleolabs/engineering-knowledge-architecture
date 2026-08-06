# View Fixture

Convention document used by view integration tests.

The fixture is a conformant repository (asserted by `TestFixtureConforms`)
with two execution containers in namespace `eka-view-fixture`:

- `ctr:wave-0` — completed; holds `sto:legacy` via `tkt-sto-legacy`.
- `ctr:wave-1` — active; holds `sto:alpha`, `sto:beta`, `ts:gamma`,
  `bug:delta`, `ch:epsilon` via their tickets.

Tickets (all `tkt-`, under `docs/operating/projections/`):

| Ticket | derives-from | Fixture role |
|---|---|---|
| tkt-sto-alpha | ctr:wave-1, sto:alpha | regular ticket |
| tkt-sto-alpha-dup | ctr:wave-1, sto:alpha | second ticket referencing the already-covered `sto:alpha` — exercises dedup-by-identity in `WorkItemsForContainer` (the work item must appear once) |
| tkt-sto-beta | ctr:wave-1, sto:beta | regular ticket |
| tkt-sto-beta-multi | ctr:wave-1, sto:beta, ts:gamma | TWO work item references — exercises first-resolvable-wins in `ticketTargets` (`sto:beta` must win deterministically) |
| tkt-ts-gamma | ctr:wave-1, ts:gamma | regular ticket |
| tkt-bug-delta | ctr:wave-1, bug:delta | regular ticket |
| tkt-ch-epsilon | ctr:wave-1, ch:epsilon | regular ticket |
| tkt-unresolved | ctr:wave-1 | no work item reference — projects "unresolved" |
| tkt-sto-legacy | ctr:wave-0, sto:legacy | ticket of the completed container |
