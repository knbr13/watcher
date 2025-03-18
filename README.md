# Watcher 👁️

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/)
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
go install github.com/knbr13/watcher@latest
```

### Prebuilt Binaries
Download from [Releases](https://github.com/knbr13/watcher/releases)

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

## Configuration Guide 📋

### Example Config
```yaml
# Global hooks
on_success: "echo 'All systems go! 🚀'"
on_failure: "curl -X POST https://api.status.io/down"

write:
  - pattern: "src/**/*.js"
    commands:
      - "npm run lint"
      - "npm run build"
    sequential: true  # Run commands in order
    timeout: 1m       # Fail if build takes >1 minute

create:
  - pattern: "*.{jpg,png}"
    commands: ["convert $FILE -resize 50% resized/$FILE_BASE"]
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
      - "clamscan $FILE --move=/quarantine"
      - "curl -X POST http://localhost:8080/alert -d 'Infected: $FILE'"
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
--file, -f     Configuration file (required)
--path, -p     Directory to watch (default: current)
--recursive, -r Watch directories recursively
```

## Acknowledgements 💛

Built with these awesome libraries:
- [fsnotify](https://github.com/fsnotify/fsnotify) - File system notifications
- [go-arg](https://github.com/alexflint/go-arg) - CLI argument parsing
- [doublestar](https://github.com/bmatcuk/doublestar) - Glob pattern matching

---

**Watcher** © 2025 - MIT License | Crafted with ❤️ by [knbr13](https://github.com/knbr13)
