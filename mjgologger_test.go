package mjgologger

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestSetupAndStop(t *testing.T) {
	testFile := "test_setup.log"
	defer os.Remove(testFile)

	Setup(testFile)

	if myLogger == nil {
		t.Error("Setup() failed to initialize logger")
	}

	if file == nil {
		t.Error("Setup() failed to open file")
	}

	Stop()

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Setup() did not create log file")
	}
}

func TestInfo(t *testing.T) {
	testFile := "test_info.log"
	defer os.Remove(testFile)

	Setup(testFile)
	defer Stop()

	Info("Test info message")
	Info("Test with args: %s %d", "hello", 42)

	// Read and verify log content
	file.Sync() // Ensure data is written
	content := readLogFile(t, testFile)

	if !strings.Contains(content, "[INFO]") {
		t.Error("Info() did not write [INFO] prefix")
	}
	if !strings.Contains(content, "Test info message") {
		t.Error("Info() did not write message")
	}
	if !strings.Contains(content, "hello 42") {
		t.Error("Info() did not format args correctly")
	}
}

func TestDebug(t *testing.T) {
	testFile := "test_debug.log"
	defer os.Remove(testFile)

	Setup(testFile)
	defer Stop()

	Debug("Test debug message")
	Debug("Debug with number: %d", 123)

	file.Sync()
	content := readLogFile(t, testFile)

	if !strings.Contains(content, "[DEBUG]") {
		t.Error("Debug() did not write [DEBUG] prefix")
	}
	if !strings.Contains(content, "Test debug message") {
		t.Error("Debug() did not write message")
	}
}

func TestWarn(t *testing.T) {
	testFile := "test_warn.log"
	defer os.Remove(testFile)

	Setup(testFile)
	defer Stop()

	Warn("Test warning message")
	Warn("Warning with value: %v", true)

	file.Sync()
	content := readLogFile(t, testFile)

	if !strings.Contains(content, "[WARN]") {
		t.Error("Warn() did not write [WARN] prefix")
	}
	if !strings.Contains(content, "Test warning message") {
		t.Error("Warn() did not write message")
	}
}

func TestError(t *testing.T) {
	testFile := "test_error.log"
	defer os.Remove(testFile)

	Setup(testFile)
	defer Stop()

	Error("Test error message")
	Error("Error with code: %d", 500)

	file.Sync()
	content := readLogFile(t, testFile)

	if !strings.Contains(content, "[ERROR]") {
		t.Error("Error() did not write [ERROR] prefix")
	}
	if !strings.Contains(content, "Test error message") {
		t.Error("Error() did not write message")
	}
}

func TestCallerInfo(t *testing.T) {
	testFile := "test_caller.log"
	defer os.Remove(testFile)

	Setup(testFile)
	defer Stop()

	Info("Test caller information")

	file.Sync()
	content := readLogFile(t, testFile)

	// Verify that caller information is included (file path and line number)
	if !strings.Contains(content, "mjgologger_test.go") {
		t.Error("Log did not include caller file information")
	}
	if !strings.Contains(content, "[") || !strings.Contains(content, "]") {
		t.Error("Log did not include proper formatting for caller info")
	}
}

func TestMultipleLogLevels(t *testing.T) {
	testFile := "test_multiple.log"
	defer os.Remove(testFile)

	Setup(testFile)
	defer Stop()

	Info("Info message")
	Debug("Debug message")
	Warn("Warn message")
	Error("Error message")

	file.Sync()
	content := readLogFile(t, testFile)

	// Verify all log levels are present
	if !strings.Contains(content, "[INFO]") {
		t.Error("Missing [INFO] level")
	}
	if !strings.Contains(content, "[DEBUG]") {
		t.Error("Missing [DEBUG] level")
	}
	if !strings.Contains(content, "[WARN]") {
		t.Error("Missing [WARN] level")
	}
	if !strings.Contains(content, "[ERROR]") {
		t.Error("Missing [ERROR] level")
	}
}

func TestSetupWithInvalidPath(t *testing.T) {
	// Try to create a file in an invalid directory
	invalidPath := "Z:\\nonexistent\\directory\\test.log"

	Setup(invalidPath)

	// The function should handle the error gracefully
	// and myLogger should be nil or not panic on subsequent calls
	if myLogger != nil {
		t.Log("Setup with invalid path still created logger (may be expected behavior)")
	}
}

// Helper function to read log file content
func readLogFile(t *testing.T, filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	defer f.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	return content.String()
}
