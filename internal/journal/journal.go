package journal

import (
	"bufio"
	"os/exec"

	"ssh-audit-exporter/logger"
)

func TailSSHJournal(parse func(string)) error {

	logger.Info("Starting journalctl stream")

	cmd := exec.Command(
		"journalctl",
		"-f",
		"-o",
		"cat",
		"_SYSTEMD_UNIT=ssh.service",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			parse(line)
		}
	}

	return scanner.Err()
}