package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/VandiKond/vanerrors/vanstack"
)

// The file names
const (
	INFO_a_WARN   string = "logs/info.txt"
	ERROR_a_FATAL string = "logs/error.txt"
)

// The errors
const ()

// The time format
const (
	FORMAT string = "02.01.06 3:04:05 "
)

// Log levels
type LogLevel int

// Log level
const (
	INFO LogLevel = iota
	WARN
	ERROR
	FATAL
)

// Standard log level
var StringLogLevel = map[LogLevel]string{
	INFO:  ">>info<<:",
	WARN:  "!!warn!!:",
	ERROR: "**error**:",
	FATAL: "##fatal##:",
}

// The logger
type Logger struct {
	wInfo    io.Writer
	wWarn    io.Writer
	wError   io.Writer
	wFatal   io.Writer
	levelMap map[LogLevel]string
}

// Creates a logger with pairs
func New() *Logger {
	wIaW, err := os.OpenFile(INFO_a_WARN, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	wEaF, err := os.OpenFile(ERROR_a_FATAL, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	logger := Logger{
		wInfo:    wIaW,
		wWarn:    wIaW,
		wError:   wEaF,
		wFatal:   wEaF,
		levelMap: StringLogLevel,
	}
	fmt.Fprint(wIaW, "\n\n")
	fmt.Fprint(wEaF, "\n\n")
	return &logger
}

func writeln(w io.Writer, prefix string, a []any) {
	fmt.Fprintln(w, append([]any{prefix, time.Now().Format(FORMAT)}, a...)...)
}

func writef(w io.Writer, prefix string, format string, a []any) {
	format = "%s %s " + format + "\n"
	fmt.Fprintf(w, format, append([]any{prefix, time.Now().Format(FORMAT)}, a...)...)
}

// Prints a line
func (l *Logger) Println(a ...any) {
	writeln(l.wInfo, l.levelMap[INFO], a)
}

// Prints a formatted line
func (l *Logger) Printf(format string, a ...any) {
	writef(l.wInfo, l.levelMap[INFO], format, a)
}

// Prints a warn line
func (l *Logger) Warnln(a ...any) {
	writeln(l.wWarn, l.levelMap[WARN], a)
}

// Prints a warn formatted line
func (l *Logger) Warnf(format string, a ...any) {
	writef(l.wWarn, l.levelMap[WARN], format, a)
}

// Prints a error line
func (l *Logger) Errorln(a ...any) {
	writeln(l.wError, l.levelMap[ERROR], a)
}

// Prints a error formatted line
func (l *Logger) Errorf(format string, a ...any) {
	writef(l.wError, l.levelMap[ERROR], format, a)
}

// Prints a fatal line and exit
func (l *Logger) Fatalln(a ...any) {
	writeln(l.wFatal, l.levelMap[FATAL], a)
	stack := vanstack.NewStack()
	stack.Fill("", 20)
	fmt.Fprintln(os.Stderr, stack)
	os.Exit(http.StatusTeapot)
}

// Prints a fatal formatted line and exit
func (l *Logger) Fatalf(format string, a ...any) {
	writef(l.wFatal, l.levelMap[FATAL], format, a)
	stack := vanstack.NewStack()
	stack.Fill("", 20)
	fmt.Fprintln(os.Stderr, stack)
	os.Exit(http.StatusTeapot)
}
