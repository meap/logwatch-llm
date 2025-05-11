package config

import (
	"github.com/sirupsen/logrus"
)

type PlainFormatter struct{}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func SetupLogging() {
	logrus.SetFormatter(new(PlainFormatter))
}
