package journal

import (
	"github.com/coreos/go-systemd/v22/sdjournal"
)

func TailSSHJournal(parse func(string)) error {
	j, err := sdjournal.NewJournal()
	if err != nil {
		return err
	}
	defer j.Close()

	_ = j.SeekTail()

	for {
		n, err := j.Next()
		if err != nil {
			return err
		}
		if n == 0 {
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