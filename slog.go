package package_log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

//var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

//
//func GetLogger() *Logger {
//	return &Logger{e}
//}

//func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
//	return &Logger{l.WithField(k, v)}
//}

func GetLogger(filePath string, fileName string) *Logger {

	l := logrus.New()
	l.SetReportCaller(true)

	l.SetFormatter(&logrus.JSONFormatter{})

	//l.Formatter = &logrus.TextFormatter{
	//	CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
	//		filename := path.Base(frame.File)
	//		return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
	//	},
	//	DisableColors: false,
	//	FullTimestamp: true,
	//}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, 0777)
		if err != nil {
			panic(err)
		}
	}

	allFile, err := os.OpenFile(filePath+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	return &Logger{logrus.NewEntry(l)}
}
