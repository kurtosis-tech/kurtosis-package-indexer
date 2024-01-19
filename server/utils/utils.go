package utils

import (
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
)

func GetFromEnvVar(
	key string,
	subject string,
) (string, error) {
	value := os.Getenv(key)
	if len(value) < 1 {
		return "", stacktrace.NewError("No '%s' env var was found. Must be provided as env var %s", subject, key)
	}
	logrus.Debugf("Successfully loaded env var '%s'", subject)
	return value, nil
}
