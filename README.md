# MjGoLogger

A Go library for file-based logging with multiple severity levels and automatic caller information.

## Features

- Multiple log levels: Info, Debug, Warn, Error
- Automatic caller information (file and line number)
- File-based logging with timestamps
- Simple setup and teardown

## Installation

```bash
go get github.com/mjbmarques/MjGoLogger
```

## Usage

### Basic Example

```go
package main

import "github.com/mjbmarques/MjGoLogger"

func main() {
    // Setup logger with a log file
    mjgologger.Setup("app.log")
    defer mjgologger.Stop()

    // Log messages at different levels
    mjgologger.Info("Application started")
    mjgologger.Debug("Debug information: %s", "some value")
    mjgologger.Warn("This is a warning")
    mjgologger.Error("An error occurred: %v", err)
}
```

## API Reference

### Setup

Initializes the logger with the specified log file. Must be called before using any logging functions.

```go
func Setup(fileName string)
```

**Parameters:**
- `fileName`: Path to the log file to create/write to

**Example:**
```go
mjgologger.Setup("logs/app.log")
```

### Stop

Closes the log file. Should be called when logging is complete (typically with `defer`).

```go
func Stop()
```

**Example:**
```go
mjgologger.Setup("app.log")
defer mjgologger.Stop()
```

### Info

Logs an informational message with [INFO] severity level.

```go
func Info(msg string, args ...any)
```

**Parameters:**
- `msg`: Message format string
- `args`: Optional arguments for string formatting

**Example:**
```go
mjgologger.Info("User %s logged in", username)
```

### Debug

Logs a debug message with [DEBUG] severity level.

```go
func Debug(msg string, args ...any)
```

**Parameters:**
- `msg`: Message format string
- `args`: Optional arguments for string formatting

**Example:**
```go
mjgologger.Debug("Variable value: %v", someVar)
```

### Warn

Logs a warning message with [WARN] severity level.

```go
func Warn(msg string, args ...any)
```

**Parameters:**
- `msg`: Message format string
- `args`: Optional arguments for string formatting

**Example:**
```go
mjgologger.Warn("Connection timeout, retrying...")
```

### Error

Logs an error message with [ERROR] severity level.

```go
func Error(msg string, args ...any)
```

**Parameters:**
- `msg`: Message format string
- `args`: Optional arguments for string formatting

**Example:**
```go
mjgologger.Error("Failed to connect: %v", err)
```

## Log Format

Each log entry includes:
- Timestamp (date and time)
- Severity level ([INFO], [DEBUG], [WARN], [ERROR])
- Caller information (file path and line number)
- Log message

**Example log output:**
```
2025/12/03 15:04:05 [INFO][/path/to/file.go:42]: Application started
2025/12/03 15:04:06 [DEBUG][/path/to/file.go:43]: Processing request ID: 12345
2025/12/03 15:04:07 [WARN][/path/to/file.go:44]: Cache miss for key: user_data
2025/12/03 15:04:08 [ERROR][/path/to/file.go:45]: Database connection failed
```

## License

See LICENSE file for details.

