package metrics

import "github.com/sirupsen/logrus"

// doNothingReporter will do noting when is called
// it was created to avoid calling Snowflake from the CI builds
type doNothingReporter struct{}

func (reporter *doNothingReporter) Schedule(_ bool) error {
	logrus.Debugf("doNothingReporter has been called, nothing to schedulle")
	return nil
}
