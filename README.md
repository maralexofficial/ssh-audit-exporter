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