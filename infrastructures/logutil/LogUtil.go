package logutil

import (
    "errors"
    "io"
    "log"
    "os"
)

const logFile = "./logs.txt"

func Setup() {
    // Delete log
    err := os.Remove(logFile)
    if err != nil {
        if !errors.Is(err, os.ErrNotExist) {
            Fatal("Failed to delete the log file. Error was %v", err)
        }
    }
    // If the file doesn't exist, create it or append to the file
    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        Fatal("Failed to open log file. Error was %v", err)
    }

    mw := io.MultiWriter(os.Stdout, file)
    log.SetOutput(mw)
}

func Info(msg string) {
    log.Printf("INFO: %v", msg)
}

func Warn(msg string) {
    log.Printf("WARN: %v", msg)
}

func Error(msg string, a ...interface{}) {
    m := "ERROR: " + msg
    log.Printf(m, a...)
}

func Fatal(msg string, a ...interface{}) {
    m := "ERROR: " + msg
    log.Fatalf(m, a...)
}
