package utils

import (
	"github.com/sirupsen/logrus"
	"hw7/conf"
	"log/syslog"
	"os"
)

// NewLogger - Создаёт новый логгер
func NewLogger(conf *conf.Config) (*logrus.Logger, error) {
	lg := logrus.New()

	level, err := logrus.ParseLevel(conf.Logger.Level)
	if err != nil {
		return nil, err
	}
	lg.SetLevel(level)

	lg.SetReportCaller(conf.Logger.ReportCaller)
	lg.SetFormatter(&logrus.TextFormatter{})

	if conf.Logger.Syslog {
		sysw, err := syslog.New(syslog.LOG_DEBUG, "serv")
		if err != nil {
			return nil, err
		}
		lg.SetOutput(sysw)
	} else {
		if conf.Logger.Output != "" {
			f, err := os.Create(conf.Logger.Output)
			if err != nil {
				return nil, err
			}
			lg.SetOutput(f)
		}
	}

	return lg, nil
}
