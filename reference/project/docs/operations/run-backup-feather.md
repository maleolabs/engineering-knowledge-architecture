---
namespace: feather
type: run
id: backup-feather
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: operations
author: Amira Khan
created: 2026-06-26
updated: 2026-07-02
supersedes: []
derives-from: []
depends-on:
  - arc:feather-system
amends: []
validates: []
change-log:
  - date: 2026-06-26
    domain: existence-state
    from: "-"
    to: active
    by: Amira Khan
  - date: 2026-06-26
    domain: content-state
    from: "-"
    to: draft
    by: Amira Khan
  - date: 2026-06-29
    domain: content-state
    from: draft
    to: review
    by: Amira Khan
  - date: 2026-07-02
    domain: content-state
    from: review
    to: approved
    by: Amira Khan
---

# Runbook — Backup Feather

## Purpose

Back up a Feather installation and restore it. The storage model (`adr:content-storage`) makes this a two-part job: the post Markdown files (source of truth) and the SQLite index (derived, but lossy to rebuild if the FTS state or revisions matter). Both must be captured; the files are the recovery-critical part.

## Content

**Backup (daily cron, `~/feather/backup.sh`):**

1. **Consistent SQLite snapshot** — use the online backup API, never a plain file copy (writes are atomic but concurrent):
   ```sh
   sqlite3 ~/feather/feather.db ".backup '/tmp/feather.db.bak'"
   ```
2. **Copy the post files:**
   ```sh
   tar -C ~/feather -czf /tmp/posts-$(date +%F).tar.gz content/posts/
   ```
3. **Package and ship off-host:**
   ```sh
   tar -czf /tmp/feather-backup-$(date +%F).tar.gz /tmp/feather.db.bak /tmp/posts-*.tar.gz
   # copy to the backup bucket (rclone / rsync to the off-site target)
   ```
4. **Retention:** keep 30 daily backups + the last 12 monthly; verify the newest file's checksum after every upload.

**Restore (recovery, worst case: fresh host):**

1. Install the release tag (steps 1–2 of `run:deploy-feather`).
2. Restore the posts first (the source of truth): extract `content/posts/` into `~/feather/content/`.
3. Restore the database from the backup; if the DB is missing or older than the posts, run `./feather reindex` — the index is derived and rebuilds from the files.
4. Start the service, then verify: post count matches pre-backup, search finds a known post, revisions intact.

**Expected results:** daily off-host backups, restorable to a fresh host within 30 minutes; index loss is never data loss (files are the source of truth).
