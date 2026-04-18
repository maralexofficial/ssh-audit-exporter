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

	_ = j.SeekTail()

	_ = j.AddMatch("_SYSTEMD_UNIT=ssh.service")

	for {
		n, err := j.Next()
		if err != nil {
			return err
		}

		if n == 0 {
			j.Wait(sdjournal.IndefiniteWait)
			continue
		}

		entry, err := j.GetEntry()
		if err != nil {
			continue
		}

		msg := entry.Fields["MESSAGE"]
		if msg != "" {
			parse(msg)
		}
	}
}