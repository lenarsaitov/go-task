package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: false,
	}

	//err := os.MkdirAll("logs", 0755)
	//if err != nil || os.IsExist(err){
	//	l.Panic("Can't create log dir. no configured logging to files")
	//}

	//_, err = os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	//if err != nil{
	//	l.Panicf("Cant open log file: %s", err)
	//}
	//l.SetOutput(allFile)

	l.SetLevel(logrus.TraceLevel)
	e = logrus.NewEntry(l)
}
