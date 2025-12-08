package mjgologger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var myLogger *log.Logger
var file *os.File

func Setup(fileName string) error {

	// Take care of log rotation by renaming existing file
	if fileExists(fileName) {
		renameFile(fileName)
	}

	// Create new log file
	var err error
	file, err = os.Create(fileName)
	if err != nil {
		fmt.Println("ERROR: failed to create file with filename: ", fileName)
		return err
	}
	myLogger = log.New(file, "", log.LstdFlags)
	return nil
}

func Stop() error {
	err := file.Close()
	if err != nil {
		fmt.Println("Error: failed to close file with filename: " + file.Name())
		return err
	}
	return nil
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

func renameFile(oldName string) {
	// Append timestamp to old file name with format YYYYMMDD_HHMMSS.ms, based on time lib's format.go reference.
	timestamp := time.Now().Format("20060102_150405.00")
	newFileName := oldName + "." + timestamp
	err := os.Rename(oldName, newFileName)
	if err != nil {
		fmt.Println("ERROR: failed to rename existing file: ", oldName, "to", newFileName)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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
