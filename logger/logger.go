package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/sirupsen/logrus"
)

var Log *logger
var logT logInterface // for text formatter
var logJ logInterface // for json formatter
var once sync.Once

// logger empty struct for debugging purposes
type logger struct{}

type logTextFormatter struct {
	logrus.TextFormatter
}

type logInterface interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Println(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
}

func init() {
	once.Do(func() {
		logT = initTextFormatter()
		logJ = initJsonFormatter()
	})
}

func initTextFormatter() *logrus.Logger {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logTextFormatter{
		logrus.TextFormatter{
			ForceColors:            true,
			FullTimestamp:          true,
			TimestampFormat:        time.RFC3339,
			DisableLevelTruncation: false,
		},
	})

	return log
}

func initJsonFormatter() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logJsonFormatter{
		logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   time.RFC3339,
		},
	})

	return log
}

func (l *logTextFormatter) Format(e *logrus.Entry) ([]byte, error) {
	var (
		levelColor int
	)

	switch e.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 35 // purple
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}

	return []byte(
		fmt.Sprintf(
			"%s \x1b[%dm%s\x1b[0m - %s\n",
			e.Time.Format(l.TimestampFormat),
			levelColor,
			strings.ToUpper(e.Level.String()),
			e.Message,
		),
	), nil
}

func (l *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	var (
		lms  []LogMessage
		file string
	)

	if env.GetBool("DEBUG", true) {
		logT.Errorf(format, args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.ErrorLevel.String(),
		Message: fmt.Sprintf(format, args...),
	})

	value.Set(_LogMessage, lms)
}

func (l *logger) Error(ctx context.Context, args ...interface{}) {
	var (
		lms  []LogMessage
		file string
	)

	if env.GetBool("DEBUG", true) {
		logT.Error(args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.ErrorLevel.String(),
		Message: fmt.Sprint(args...),
	})

	value.Set(_LogMessage, lms)
}

func (l *logger) ErrorWithFilename(ctx context.Context, file string, args ...interface{}) {
	var (
		lms []LogMessage
	)

	if env.GetBool("DEBUG", true) {
		logT.Error(args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.ErrorLevel.String(),
		Message: fmt.Sprint(args...),
	})

	value.Set(_LogMessage, lms)
}

func (l *logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	var (
		lms  []LogMessage
		file string
	)

	if env.GetBool("DEBUG", true) {
		logT.Debugf(format, args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.DebugLevel.String(),
		Message: fmt.Sprintf(format, args...),
	})

	if env.GetBool("DEBUG", true) {
		value.Set(_LogMessage, lms)
	}
}

func (l *logger) Debug(ctx context.Context, args ...interface{}) {
	var (
		lms  []LogMessage
		file string
	)

	if env.GetBool("DEBUG", true) {
		logT.Debug(args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.DebugLevel.String(),
		Message: fmt.Sprint(args...),
	})

	if env.GetBool("DEBUG", true) {
		value.Set(_LogMessage, lms)
	}
}

func (l *logger) DebugWithFilename(ctx context.Context, file string, args ...interface{}) {
	var (
		lms []LogMessage
	)

	if env.GetBool("DEBUG", true) {
		logT.Debug(args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.DebugLevel.String(),
		Message: fmt.Sprint(args...),
	})

	if env.GetBool("DEBUG", true) {
		value.Set(_LogMessage, lms)
	}
}

func (l *logger) Printf(ctx context.Context, format string, args ...interface{}) {
	var (
		lms  []LogMessage
		file string
	)

	if env.GetBool("DEBUG", true) {
		logT.Println(fmt.Sprintf(format, args...))
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.InfoLevel.String(),
		Message: fmt.Sprintf(format, args...),
	})

	value.Set(_LogMessage, lms)
}

func (l *logger) Print(ctx context.Context, args ...interface{}) {
	var (
		lms  []LogMessage
		file string
	)

	if env.GetBool("DEBUG", true) {
		logT.Println(args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get file and line
	_, fileName, line, _ := runtime.Caller(1)
	// default file
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.InfoLevel.String(),
		Message: fmt.Sprint(args...),
	})

	value.Set(_LogMessage, lms)
}

func (l *logger) PrintWithFilename(ctx context.Context, file string, args ...interface{}) {
	var (
		lms []LogMessage
	)

	if env.GetBool("DEBUG", true) {
		logT.Println(args...)
	}

	// get current log from context
	value, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := value.LoadAndDelete(_LogMessage)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// append data to lms
	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.InfoLevel.String(),
		Message: fmt.Sprint(args...),
	})

	value.Set(_LogMessage, lms)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	logT.Fatalf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	logT.Fatal(args...)
}

type logJsonFormatter struct {
	logrus.JSONFormatter
}

// DB start logging
func DB(types DatabaseType, query string, args ...interface{}) Database {
	startTime := time.Now()

	// Convert the elements to strings and format them
	var strElements []string
	for _, v := range args {
		val := fmt.Sprintf("%v", v)
		if len(val) > 1000 {
			val = "char argument too much"
		}

		strElements = append(strElements, val)
	}

	return Database{
		Type:      types,
		Query:     query,
		Arguments: strElements,
		startTime: startTime,
	}
}

func RedBold(str interface{}) {
	fmt.Printf("\x1b[31;1m%v\x1b[0m\n", str)
}

func RedItalic(str interface{}) {
	fmt.Printf("\x1b[31;3m%v\x1b[0m\n", str)
}

func Red(str interface{}) {
	fmt.Printf("\x1b[31;5m%v\x1b[0m\n", str)
}

func GreenBold(str interface{}) {
	fmt.Printf("\x1b[32;1m%v\x1b[0m\n", str)
}

func GreenItalic(str interface{}) {
	fmt.Printf("\x1b[32;3m%v\x1b[0m\n", str)
}

func Green(str interface{}) {
	fmt.Printf("\x1b[32;5m%v\x1b[0m\n", str)
}

func YellowBold(str interface{}) {
	fmt.Printf("\x1b[33;1m%v\x1b[0m\n", str)
}

func YellowItalic(str interface{}) {
	fmt.Printf("\x1b[33;3m%v\x1b[0m\n", str)
}

func Yellow(str interface{}) {
	fmt.Printf("\x1b[33;5m%v\x1b[0m\n", str)
}

func BlueBold(str interface{}) {
	fmt.Printf("\x1b[34;1m%v\x1b[0m\n", str)
}

func BlueItalic(str interface{}) {
	fmt.Printf("\x1b[34;3m%v\x1b[0m\n", str)
}

func Blue(str interface{}) {
	fmt.Printf("\x1b[34;5m%v\x1b[0m\n", str)
}

func PurpleBold(str interface{}) {
	fmt.Printf("\x1b[35;1m%v\x1b[0m\n", str)
}

func PurpleItalic(str interface{}) {
	fmt.Printf("\x1b[35;3m%v\x1b[0m\n", str)
}

func Purple(str interface{}) {
	fmt.Printf("\x1b[35;5m%v\x1b[0m\n", str)
}

func CyanBold(str interface{}) {
	fmt.Printf("\x1b[36;1m%v\x1b[0m\n", str)
}

func CyanItalic(str interface{}) {
	fmt.Printf("\x1b[36;3m%v\x1b[0m\n", str)
}

func Cyan(str interface{}) {
	fmt.Printf("\x1b[36;5m%v\x1b[0m\n", str)
}

func SetRequestId(ctx context.Context, requestId string) {
	value, ok := extract(ctx)
	if !ok {
		return
	}

	value.Set(_RequestId, requestId)
}

func GetRequestId(ctx context.Context) string {
	value, ok := extract(ctx)
	if !ok {
		return ""
	}

	val, ok := value.Load(_RequestId)
	if ok {
		return val.(string)
	}

	return uuid.NewString()
}

func GetUsername(ctx context.Context) string {
	value, ok := extract(ctx)
	if !ok {
		return ""
	}

	val, ok := value.Load(_Username)
	if ok {
		return val.(string)
	}

	return ""
}
