package logger

import (
	"fmt"
	"log"
	"runtime"
)

func Info(msg string) {
	log.Println(msg)
}

func Error(err any) {
	_, file, line, ok := runtime.Caller(1)

	if !ok {
		log.Println(err)
		return
	}

	log.Printf("%s:%d: %s\n", file, line, err)
}

func Errorf(format string, values ...any) {
	_, file, line, ok := runtime.Caller(1)
	msg := fmt.Sprintf(format, values...)

	if ok {

		log.Printf("%s:%d: %s\n", file, line, msg)
		return
	}

	log.Println(msg)
}
