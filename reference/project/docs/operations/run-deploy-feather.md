---
namespace: feather
type: run
id: deploy-feather
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: operations
author: Amira Khan
created: 2026-06-25
updated: 2026-07-01
supersedes: []
derives-from: []
depends-on:
  - arc:feather-system
amends: []
validates: []
change-log:
  - date: 2026-06-25
    domain: existence-state
    from: "-"
    to: active
    by: Amira Khan
  - date: 2026-06-25
    domain: content-state
    from: "-"
    to: draft
    by: Amira Khan
  - date: 2026-06-28
    domain: content-state
    from: draft
    to: review
    by: Amira Khan
  - date: 2026-07-01
    domain: content-state
    from: review
    to: approved
    by: Amira Khan
---

# Runbook — Deploy Feather

## Purpose

Deploy a new Feather release to the production host behind Caddy (`dec:reverse-proxy`): build the binary, run migrations, restart the service, and verify. Prerequisites: SSH access to the host, the release artifact tag, and the deploy standard checklist from `std:definition-of-done`.

## Content

**Steps (run on the production host, from the deploy directory `~/feather`):**

1. **Fetch and verify the release:**
   ```sh
   git -C ~/feather checkout <release-tag> && git -C ~/feather pull --ff-only
   ```
   Verify the tag exists and the working tree is clean.

2. **Build the binary (Go 1.24+, `ch:update-go-version`):**
   ```sh
   cd ~/feather && CGO_ENABLED=0 go build -trimpath -o /tmp/feather.new ./cmd/feather
   /tmp/feather.new version   # sanity: prints build + EKA standard version
   ```

3. **Migrate the database (idempotent):**
   ```sh
   /tmp/feather.new migrate --db ~/feather/feather.db
   ```
   Migrations run forward-only; a failed migration aborts the deploy (do not proceed).

4. **Swap and restart:**
   ```sh
   mv /tmp/feather.new ~/feather/feather
   systemctl restart feather
   systemctl --no-pager status feather   # expect: active (running)
   ```

5. **Verify behind Caddy (TLS terminates there; certificates auto-renew):**
   ```sh
   curl -fsSI https://feather.example.com/           # expect HTTP 200/304
   curl -fsS https://feather.example.com/feed.xml    # if distribution epic shipped
   journalctl -u feather -n 20 --no-pager            # no panics in the tail
   ```

**Rollback:** `git checkout <previous-tag>` + rebuild + restart (same steps 2–4). The DB is forward-migrated; roll back only if no migration ran, otherwise restore from `run:backup-feather` and redeploy.

**Expected results:** new version serving, health checks green, zero downtime beyond the restart window (< 5 s).
