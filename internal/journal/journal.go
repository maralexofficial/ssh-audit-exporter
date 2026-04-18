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

	_ = j.SeekHead()
	// _ = j.SeekTail()

	// DEBUG: KEIN FILTER
	// _ = j.AddMatch("_COMM=sshd")

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

		logger.Info("---- JOURNAL ENTRY ----")

		for k, v := range entry.Fields {
			logger.Info(k + "=" + v)
		}

		msg := entry.Fields["MESSAGE"]
		if msg != "" {
			logger.Info("MESSAGE=" + msg)
			parse(msg)
		}
	}
}