# Linux Log Analysis — Command Cheatsheet

All commands below are read-only. Never modify, rotate, truncate, or delete
logs; never stop or restart services/containers.

> Be careful with secrets/PII in logs. Reference by key/context, do not
> reproduce sensitive values.

## 1. journald (systemd)

```bash
# Errors and above, last hour
journalctl -p err..alert --since "1 hour ago"

# Logs for a specific service
journalctl -u nginx.service --since today

# Since last boot, kernel messages
journalctl -k -b

# Precise time window
journalctl --since "2026-09-26 08:00" --until "2026-09-26 09:30"

# By PID or executable
journalctl _PID=1234
journalctl /usr/sbin/sshd

# JSON output for post-processing
journalctl -u myapp -o json --since "30 min ago"
```

## 2. /var/log files

```bash
# Authentication / denied access (Debian: auth.log | RHEL: secure)
grep -Ei "failed|invalid|denied|refused" /var/log/auth.log

# Kernel: OOM, segfault, panic
grep -Ei "oom|out of memory|segfault|panic|call trace" /var/log/kern.log

# Compressed / rotated logs
zgrep -Ei "error|fatal" /var/log/syslog.*.gz
zcat /var/log/nginx/error.log.2.gz | tail -n 100

# Nginx / Apache
tail -n 200 /var/log/nginx/error.log
grep -E " 5[0-9]{2} " /var/log/nginx/access.log | awk '{print $9}' | sort | uniq -c | sort -rn

# Top IPs in access log
awk '{print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head
```

## 3. Common search patterns

```bash
# Common errors in any log
grep -Ein "error|fatal|panic|timeout|refused|denied|exception|traceback" <file>

# Count occurrences by type
grep -Eio "error|warn|fatal" app.log | sort | uniq -c

# Context around an error (3 lines before/after)
grep -n -C3 "OutOfMemory" app.log
```

## 4. Containers — Docker / Podman

```bash
# Recent logs with timestamps and a time window
docker logs --since 30m --timestamps <container>
docker logs --tail 200 <container>
podman logs --since 1h -t <container>

# Compose
docker compose logs --since 15m --tail 100 <service>

# Container state: restart count, OOMKilled, health
docker inspect -f '{{.RestartCount}} {{.State.OOMKilled}} {{.State.Health.Status}}' <container>

# Runtime events (restarts, kills, health)
docker events --since 1h --filter event=die --filter event=oom
podman events --since 1h
```

## 5. Kubernetes

```bash
# Current and previous container logs (after a crash)
kubectl logs <pod> -c <container> --since=15m
kubectl logs <pod> -p          # previous container

# Events and restart / OOMKilled reasons
kubectl describe pod <pod>
kubectl get events --sort-by=.lastTimestamp

# Node level (containerd / CRI-O)
crictl ps -a
crictl logs <container-id>
```

## 6. lnav — multi-file navigation and correlation

```bash
# Open multiple logs (auto format detection + unified timeline)
lnav /var/log/syslog /var/log/nginx/error.log

# SQL query over structured logs (inside lnav, prompt `;`)
#   ;SELECT log_time, log_level, log_body FROM logline WHERE log_level >= 'error';

# Batch (non-interactive) mode
lnav -n -c ';SELECT count(*), log_level FROM logline GROUP BY log_level' /var/log/syslog
```

## 7. logcli / Loki (LogQL)

```bash
# Discover labels and streams
logcli labels
logcli series '{job="nginx"}'

# Query over a time window
logcli query '{job="myapp"} |= "ERROR"' --since=1h --limit=200

# Error rate per minute
logcli query 'sum(rate({job="myapp"} |= "ERROR" [1m]))' --since=30m

# Instant query
logcli instant-query 'count_over_time({job="nginx"} |= "500" [5m])'
```
