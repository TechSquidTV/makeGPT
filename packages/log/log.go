package log

import (
	"os"

	charmlog "github.com/charmbracelet/log"
)

var logger *charmlog.Logger

func init() {
	logger = charmlog.New(os.Stderr)
}

func GetLogger() *charmlog.Logger {
	return logger
}

func SetLevel(level charmlog.Level) {
	logger.SetLevel(level)
}
