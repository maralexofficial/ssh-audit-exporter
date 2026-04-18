package logfile

import (
	"os"
)

var defaultLogPaths = []string{
	"/var/log/auth.log",
	"/var/log/secure",
}

func GetLogFile(cliPath string) (string, error) {

	if cliPath != "" {
		return cliPath, nil
	}

	envPath := os.Getenv("SSH_LOG_FILE")
	if envPath != "" {
		return envPath, nil
	}

	for _, path := range defaultLogPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}