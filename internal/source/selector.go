package source

import "os"

type Type string

const (
	File    Type = "file"
	Journal Type = "journal"
)

func GetSourceType() Type {
	if os.Getenv("SSH_LOG_TYPE") == "journal" {
		return Journal
	}
	return File
}