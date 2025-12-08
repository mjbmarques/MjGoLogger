package mjgologger

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSetupAndStop(t *testing.T) {
	testFile := "test_setup.log"
	defer os.Remove(testFile)

	err := Setup(testFile)
	if err != nil {
		t.Fatalf("Setup() returned error: %v", err)
	}

	if myLogger == nil {
		t.Error("Setup() failed to initialize logger")
	}

	if file == nil {
		t.Error("Setup() failed to open file")
	}

	err = Stop()
	if err != nil {
		t.Errorf("Stop() returned error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Setup() did not create log file")
	}
}

func TestInfo(t *testing.T) {
	testFile := "test_info.log"
	defer os.Remove(testFile)

	err := Setup(testFile)
	if err != nil {
		t.Fatalf("Setup() returned error: %v", err)
	}
	defer func() {
		if err := Stop(); err != nil {
			t.Errorf("Stop() returned error: %v", err)
		}
	}()

	Info("Test info message")
	Info("Test with args: %s %d", "hello", 42)

	// Read and verify log content
	if err := file.Sync(); err != nil {
		t.Fatalf("file.Sync() returned error: %v", err)
	}
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

	err := Setup(testFile)
	if err != nil {
		t.Fatalf("Setup() returned error: %v", err)
	}
	defer func() {
		if err := Stop(); err != nil {
			t.Errorf("Stop() returned error: %v", err)
		}
	}()

	Debug("Test debug message")
	Debug("Debug with number: %d", 123)

	if err := file.Sync(); err != nil {
		t.Fatalf("file.Sync() returned error: %v", err)
	}
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

	err := Setup(testFile)
	if err != nil {
		t.Fatalf("Setup() returned error: %v", err)
	}
	defer func() {
		if err := Stop(); err != nil {
			t.Errorf("Stop() returned error: %v", err)
		}
	}()

	Warn("Test warning message")
	Warn("Warning with value: %v", true)

	if err := file.Sync(); err != nil {
		t.Fatalf("file.Sync() returned error: %v", err)
	}
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

	err := Setup(testFile)
	if err != nil {
		t.Fatalf("Setup() returned error: %v", err)
	}
	defer func() {
		if err := Stop(); err != nil {
			t.Errorf("Stop() returned error: %v", err)
		}
	}()

	Error("Test error message")
	Error("Error with code: %d", 500)

	if err := file.Sync(); err != nil {
		t.Fatalf("file.Sync() returned error: %v", err)
	}
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

	err := Setup(testFile)
	if err != nil {
		t.Fatalf("Setup() returned error: %v", err)
	}
	defer func() {
		if err := Stop(); err != nil {
			t.Errorf("Stop() returned error: %v", err)
		}
	}()

	Info("Test caller information")

	if err := file.Sync(); err != nil {
		t.Fatalf("file.Sync() returned error: %v", err)
	}
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

	err := Setup(testFile)
	if err != nil {
		t.Fatalf("Setup() returned error: %v", err)
	}
	defer func() {
		if err := Stop(); err != nil {
			t.Errorf("Stop() returned error: %v", err)
		}
	}()

	Info("Info message")
	Debug("Debug message")
	Warn("Warn message")
	Error("Error message")

	if err := file.Sync(); err != nil {
		t.Fatalf("file.Sync() returned error: %v", err)
	}
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

	err := Setup(invalidPath)

	// The function should return an error for invalid paths
	if err == nil {
		t.Error("Setup with invalid path should return an error")
	}
}

func TestSetupRenamesExistingFile(t *testing.T) {
	testFile := "test_rename.log"
	defer os.Remove(testFile)

	// Clean up any renamed files after test
	defer func() {
		matches, _ := filepath.Glob(testFile + ".*")
		for _, match := range matches {
			os.Remove(match)
		}
	}()

	// First Setup: Create initial log file
	err := Setup(testFile)
	if err != nil {
		t.Fatalf("First Setup() returned error: %v", err)
	}
	Info("First log entry")
	err = Stop()
	if err != nil {
		t.Fatalf("First Stop() returned error: %v", err)
	}

	// Verify first file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("First Setup() did not create log file")
	}

	// Read content from first file
	firstContent := readLogFile(t, testFile)
	if !strings.Contains(firstContent, "First log entry") {
		t.Error("First log file does not contain expected content")
	}

	// Wait a second to ensure timestamp will be different
	time.Sleep(1 * time.Second)

	// Second Setup: Should rename existing file and create new one
	err = Setup(testFile)
	if err != nil {
		t.Fatalf("Second Setup() returned error: %v", err)
	}
	Info("Second log entry")
	err = Stop()
	if err != nil {
		t.Fatalf("Second Stop() returned error: %v", err)
	}

	// Verify original filename still exists (as new file)
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Second Setup() did not create new log file")
	}

	// Verify the renamed file exists (with timestamp)
	matches, err := filepath.Glob(testFile + ".*")
	if err != nil {
		t.Fatalf("Failed to search for renamed files: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("Expected 1 renamed file, found %d", len(matches))
	}

	// Verify the renamed file has the timestamp format
	renamedFile := matches[0]
	if !strings.HasPrefix(renamedFile, testFile+".") {
		t.Errorf("Renamed file does not have expected prefix: %s", renamedFile)
	}

	// Verify renamed file contains old content
	renamedContent := readLogFile(t, renamedFile)
	if !strings.Contains(renamedContent, "First log entry") {
		t.Error("Renamed file does not contain original content")
	}
	if strings.Contains(renamedContent, "Second log entry") {
		t.Error("Renamed file should not contain new content")
	}

	// Verify new file contains only new content
	newContent := readLogFile(t, testFile)
	if strings.Contains(newContent, "First log entry") {
		t.Error("New file should not contain old content")
	}
	if !strings.Contains(newContent, "Second log entry") {
		t.Error("New file does not contain new content")
	}
}

func TestSetupMultipleRenames(t *testing.T) {
	testFile := "test_multiple_rename.log"
	defer os.Remove(testFile)

	// Clean up any renamed files after test
	defer func() {
		matches, _ := filepath.Glob(testFile + ".*")
		for _, match := range matches {
			os.Remove(match)
		}
	}()

	// Create and rename files multiple times
	for i := 1; i <= 3; i++ {
		err := Setup(testFile)
		if err != nil {
			t.Fatalf("Setup() iteration %d returned error: %v", i, err)
		}
		Info("Log entry %d", i)
		err = Stop()
		if err != nil {
			t.Fatalf("Stop() iteration %d returned error: %v", i, err)
		}
		time.Sleep(1 * time.Second) // Ensure different timestamps
	}

	// Verify original file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Final log file does not exist")
	}

	// Verify multiple renamed files exist
	matches, err := filepath.Glob(testFile + ".*")
	if err != nil {
		t.Fatalf("Failed to search for renamed files: %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("Expected 2 renamed files, found %d", len(matches))
	}

	// Verify final file contains only the latest entry
	finalContent := readLogFile(t, testFile)
	if !strings.Contains(finalContent, "Log entry 3") {
		t.Error("Final file does not contain latest entry")
	}
	if strings.Contains(finalContent, "Log entry 1") || strings.Contains(finalContent, "Log entry 2") {
		t.Error("Final file should not contain old entries")
	}
}

// Helper function to read log file content
func readLogFile(t *testing.T, filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("Failed to close log file: %v", err)
		}
	}()

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
