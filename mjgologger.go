package mjgologger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var myLogger *log.Logger
var file *os.File

func Setup(fileName string) {

	var err error
	file, err = os.Create(fileName)
	if err != nil {
		fmt.Println("ERROR: failed to create file with filename: ", fileName)
		return
	}
	myLogger = log.New(file, "", log.LstdFlags)
}

func Stop() {
	err := file.Close()
	if err != nil {
		fmt.Println("Error: failed to close file with filename: " + file.Name())
	}
}

func Info(msg string, args ...any) {
	logMessage("[INFO]", msg, args...)
}

func Debug(msg string, args ...any) {
	logMessage("[DEBUG]", msg, args...)
}

func Warn(msg string, args ...any) {
	logMessage("[WARN]", msg, args...)
}

func Error(msg string, args ...any) {
	logMessage("[ERROR]", msg, args...)
}

func logMessage(severity string, msg string, args ...any) {
	prefix := generatePrefix(severity, 3)
	myLogger.Println(prefix + fmt.Sprintf(msg, args...))
}

func generatePrefix(severity string, steps int) string {
	return severity + getCaller(steps) + ": "
}

func getCaller(steps int) string {
	_, file, line, _ := runtime.Caller(steps + 1)
	return fmt.Sprintf("[%s:%d]", file, line)
}
