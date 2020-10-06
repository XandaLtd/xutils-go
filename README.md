# xutils-go
Go utilities for rest API

## xlogger

Extends Uber's [Zap](go.uber.org/zap) logging

You may use these levels of logging, each level represents the log levels that are displayed
- DebugLevel
- InfoLevel
- WarningLevel
- ErrorLevel

Standardised Alchemy env vars for use
- LOG_LEVEL
- LOG_ERRORS_TO
- LOG_OUTPUT_TO

Log Output defaults to stdout
Log Errors defaults to stderr