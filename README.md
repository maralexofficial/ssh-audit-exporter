### 📡 FILE MODE

```
docker run -d \
  --name ssh-exporter \
  -p 9100:9100 \
  -v /var/log/auth.log:/var/log/auth.log:ro \
  --env-file .env \
  ssh-audit-exporter
```

### 📡 JOURNAL MODE

```
docker run -d \
  --name ssh-exporter \
  -p 9100:9100 \
  -e SSH_LOG_TYPE=journal \
  -v /run/systemd/journal/socket:/run/systemd/journal/socket:ro \
  --env-file .env \
  ssh-audit-exporter
```
### LOG TYPE

--log-type=journal
--log-type=file

SSH_LOG_TYPE=file      (default)
SSH_LOG_TYPE=journal

# SSH Logfile Paths (Linux)

## Debian / Ubuntu
/var/log/auth.log

---

## RHEL / CentOS / Fedora
/var/log/secure

---

## systemd / journald (no file-based logs)
journalctl -u ssh
journalctl -u sshd

---

# =========================
# Default SSH logs (Linux)
# =========================
```
SSH_SUCCESS_REGEX=Accepted .* for (\w+) from ([0-9.]+)
SSH_FAIL_REGEX=Failed .* for (\w+) from ([0-9.]+)
```
# =========================
# Ubuntu OpenSSH (with port info)
# =========================
```
SSH_SUCCESS_REGEX=Accepted password for (\w+) from ([0-9.]+) port \d+
SSH_FAIL_REGEX=Failed password for (\w+) from ([0-9.]+) port \d+
```
# =========================
# systemd / journald logs
# =========================
```
SSH_SUCCESS_REGEX=sshd.*Accepted .* for (\w+) from ([0-9.]+)
SSH_FAIL_REGEX=sshd.*Failed .* for (\w+) from ([0-9.]+)
```
# =========================
# Verbose / debug SSH logs
# =========================
```
SSH_SUCCESS_REGEX=.*Accepted .* user=(\w+) .* rhost=([0-9.]+)
SSH_FAIL_REGEX=.*authentication failure.*user=(\w+).*rhost=([0-9.]+)
```