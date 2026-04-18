package logfile

import (
	"github.com/hpcloud/tail"
)

func TailFile(path string, parse func(string)) error {
	t, err := tail.TailFile(path, tail.Config{
		Follow: true,
		ReOpen: true,
	})
	if err != nil {
		return err
	}

	for line := range t.Lines {
		if line != nil {
			parse(line.Text)
		}
	}

	return nil
}