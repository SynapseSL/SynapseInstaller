package logger

import (
	"log"

	"github.com/mgutz/ansi"
)

type Logger struct {
	logLevel        int
	phosphor_dbg    func(string) string
	phosphor_info   func(string) string
	phosphor_output func(string) string
	phosphor_ok     func(string) string
	phosphor_warn   func(string) string
	phosphor_err    func(string) string
	phosphor_crit   func(string) string
}

var UsedLogger = Logger{
	2, // Default to Warn and Above
	ansi.ColorFunc("white:17"),
	ansi.ColorFunc("black:white"),
	ansi.ColorFunc("white+b:232"),
	ansi.ColorFunc("black+b:22"),
	ansi.ColorFunc("black+b:190"),
	ansi.ColorFunc("white+bh:1"),
	ansi.ColorFunc("white+bhB:160"),
}

const (
	LogLevel_DEBUG    int = iota
	LogLevel_INFO     int = iota
	LogLevel_WARN     int = iota
	LogLevel_ERROR    int = iota
	LogLevel_CRITICAL int = iota
)

func (log *Logger) SetLogLevel(newLogLevel int) {
	log.logLevel = newLogLevel
}

func Debug(message string) {
	if UsedLogger.logLevel < 1 {
		log.Println(UsedLogger.phosphor_dbg("[🔍 DEBUG]") + " " + message)
	}
}

func Info(message string) {
	if UsedLogger.logLevel < 2 {
		log.Println(UsedLogger.phosphor_info("[🛈 INFO]") + " " + message)
	}
}

func Output(message string) {
	if UsedLogger.logLevel < 2 {
		log.Println(UsedLogger.phosphor_output("[🚂 OUTPUT]") + " " + message)
	}
}

func Ok(message string) {
	if UsedLogger.logLevel < 3 {
		log.Println(UsedLogger.phosphor_ok("[✅ OK]") + " " + message)
	}
}

func Warn(message string) {
	if UsedLogger.logLevel < 3 {
		log.Println(UsedLogger.phosphor_warn("[⚠️ WARN]") + " " + message)
	}
}

func Err(message string) {
	if UsedLogger.logLevel < 4 {
		log.Println(UsedLogger.phosphor_err("[🔥 ERROR]") + " " + message)
	}
}

func Critical(message string) {
	log.Fatalln(UsedLogger.phosphor_crit("[⚡ FATAL]") + " " + message)
}
