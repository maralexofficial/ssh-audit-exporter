package journal

import (
	"ssh-audit-exporter/logger"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

func TailSSHJournal(parse func(string)) error {

	logger.Info("Starting sdjournal stream")

	j, err := sdjournal.NewJournal()
	if err != nil {
		return err
	}
	defer j.Close()

	j.AddMatch("_SYSTEMD_UNIT=ssh.service")

	if err := j.SeekTail(); err != nil {
		return err
	}

	_, _ = j.Next()

	for {
		for {
			n, err := j.Next()
			if err != nil {
				return err
			}

			if n == 0 {
				break
			}

			entry, err := j.GetEntry()
			if err != nil {
				continue
			}

			msg := entry.Fields["MESSAGE"]
			if msg != "" {
				logger.Info("SSH EVENT: " + msg)
				parse(msg)
			}
		}

		j.Wait(sdjournal.IndefiniteWait)
	}
}
