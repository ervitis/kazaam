package kazaam

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Entry

func NewLogger(verbose bool) *logrus.Entry {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: false,
		FieldMap:    logrus.FieldMap{logrus.FieldKeyTime: "time"},
	})

	l.SetLevel(logrus.InfoLevel)
	if verbose {
		l.SetLevel(logrus.DebugLevel)
	}

	return logrus.NewEntry(l)
}
