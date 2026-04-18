# LOG TYPE

--log-type=journal
--log-type=file

SSH_LOG_TYPE=file      (default)
SSH_LOG_TYPE=journal

# LOG TYPE FILE

### SSH Logfile Paths (Linux)

#### Debian / Ubuntu
/var/log/auth.log
---

#### RHEL / CentOS / Fedora
/var/log/secure

---

# LOG TYPE JOURNAL

## systemd / journald (no file-based logs)
journalctl -u ssh
journalctl -u sshd

---