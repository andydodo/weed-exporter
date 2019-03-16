package logger

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

var Logger *log.Logger

func init() {
	file := initLogfile()
	Logger = log.New(file, "[wd-exporter] ", log.Ldate|log.Lmicroseconds)
}

func initLogfile() *os.File {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	// Create log dir
	if err := os.MkdirAll("/var/log/wd-exporter", os.ModePerm); err != nil {
		fmt.Printf("Create log dir failed: %s\n", err.Error())
		return nil
	}

	f, err := os.OpenFile("/var/log/wd-exporter/wd-exporter.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	return f
}
