---
namespace: feather
type: dec
id: reverse-proxy
instance-version: 1
revision: 1
content-state: accepted
existence-state: active
dimension: decisions
author: Lukas Weber
created: 2026-06-23
updated: 2026-06-30
supersedes: []
derives-from:
  - arc:feather-system
depends-on: []
amends: []
validates: []
change-log:
  - date: 2026-06-23
    domain: existence-state
    from: "-"
    to: active
    by: Lukas Weber
  - date: 2026-06-23
    domain: content-state
    from: "-"
    to: draft
    by: Lukas Weber
  - date: 2026-06-30
    domain: content-state
    from: draft
    to: accepted
    by: Lukas Weber
---

# Decision — Serve Behind Caddy

## Context

Production serving needs TLS, HTTP/2, and static file delivery in front of the Feather HTTP API. The options: run the Feather binary directly (TLS in-process via Go's crypto/tls), or put a reverse proxy in front (Caddy, nginx, or a managed load balancer).

## Decision

Run Feather **behind Caddy** in production. Caddy terminates TLS (automatic certificates), serves the static `public/` tree directly, and proxies `/api/*` to the Feather server on localhost.

## Consequences

- Certificate management disappears from the deploy runbook: Caddy renews automatically (see `run:deploy-feather`).
- Static assets are served by Caddy, so the Feather binary never holds them in memory.
- The Feather server stays bound to localhost and holds no TLS keys — a smaller attack surface and simpler local development (plain `./feather serve`).
- Cost: one extra moving part and one config file (`Caddyfile`); the proxy is now a production dependency.

## Alternatives Considered

- **TLS in-process, no proxy** — rejected: certificate renewal, OCSP, and HTTP/2 tuning would live in application code; static serving would consume the app's resources.
- **nginx** — rejected: no automatic certificate renewal; Caddy's zero-config TLS wins for a small deployment.
