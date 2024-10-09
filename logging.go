package main

import (
	"fmt"
	"os"
	"time"
)

/* Calls Log followed by os.Exit with error code. */
func Fatal(err error) {
	Log(err)
	os.Exit(1)
}

/* Printf equivalent of Fatal. */
func Fatalf(format string, args ...interface{}) {
	Logf(format, args...)
	os.Exit(1)
}

/* Appends values to log file. */
func Log(v ...interface{}) {
	f := getLogFile()
	defer f.Close()

	f.WriteString(getTimestamp() + fmt.Sprintln(v...))
}

/* Printf equivalent of Log. */
func Logf(format string, args ...interface{}) {
	Log(fmt.Sprintf(format, args...))
}

/* Opens the log file and returns a pointer to it. In case of os.PathError
 * the error message is printed to stdout.
 */
func getLogFile() *os.File {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	return logFile
}

/* Returns the timestamp for the Log function. */
func getTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05 ")
}
