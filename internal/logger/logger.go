package logger

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/gocolly/colly/v2/debug"
	"github.com/sirupsen/logrus"

	. "adc/internal/config"

	nested "github.com/antonfisher/nested-logrus-formatter"
)

type CrawlLogger struct {
	counter int32
	logger  *logrus.Logger
}

var Logger *logrus.Logger

func (l *CrawlLogger) Init() error {
	l.counter = 0
	l.logger = Logger
	l.logger.WithField("Field", "Crawler")
	return nil
}

func (l *CrawlLogger) Event(e *debug.Event) {
	i := atomic.AddInt32(&l.counter, 1)
	l.logger.Debugf("[%06d] | Clter: [%d] | Req: %d | Type: %s | %q", i, e.CollectorID, e.RequestID, e.Type, e.Values)
}

func SetupLogger() {
	// init new logger
	cfg := AdcConfig
	Logger = logrus.New()

	var level = cfg.LoggerOptions.Level
	var logPath = cfg.LoggerOptions.LogPath

	var logLevel logrus.Level
	switch level {
	case "dev":
		logLevel = logrus.DebugLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "warn":
		logLevel = logrus.WarnLevel
	default:
		logLevel = logrus.InfoLevel
	}
	// set log level
	Logger.SetLevel(logLevel)

	// set log formatter
	Logger.SetFormatter(&nested.Formatter{
		NoFieldsColors:  true,
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// set log writers
	// TODO: remove log file writer color output
	var writers []io.Writer
	file, err := os.OpenFile(
		fmt.Sprintf("%s", logPath),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	if err != nil {
		fmt.Printf("Create log file error: %v\nExpected log file path: %s\n", err, logPath)
		writers = append(writers, os.Stdout)
	} else {
		writers = append(writers, os.Stdout)
		writers = append(writers, file)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	Logger.SetOutput(fileAndStdoutWriter)
}

func init() {

}
