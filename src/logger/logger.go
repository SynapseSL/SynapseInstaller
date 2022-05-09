package logger

import (
	"log"

	"github.com/mgutz/ansi"
)

type Logger struct {
	logLevel       int
	ansiEnabled    bool
	phosphorDbg    func(string) string
	phosphorInfo   func(string) string
	phosphorOutput func(string) string
	phosphorOk     func(string) string
	phosphorWarn   func(string) string
	phosphorErr    func(string) string
	phosphorCrit   func(string) string
}

const (
	LogLevelDEBUG    int = iota
	LogLevelINFO     int = iota
	LogLevelWARN     int = iota
	LogLevelERROR    int = iota
	LogLevelCRITICAL int = iota
)

// UsedLogger is global for simplicity's sake.
// Yes, I know CI complains about it. This is fine for now.
var UsedLogger = Logger{
	LogLevelWARN,
	true,
	ansi.ColorFunc("white:17"),
	ansi.ColorFunc("black:white"),
	ansi.ColorFunc("white+b:232"),
	ansi.ColorFunc("black+b:22"),
	ansi.ColorFunc("black+b:190"),
	ansi.ColorFunc("white+bh:1"),
	ansi.ColorFunc("white+bhB:160"),
}

func (log *Logger) SetLogLevel(newLogLevel int) {
	log.logLevel = newLogLevel
}

func (log *Logger) SetAnsi(ansiEnabled bool) {
	log.ansiEnabled = ansiEnabled
}

func Debug(message string) {
	if UsedLogger.logLevel <= LogLevelDEBUG {
		if UsedLogger.ansiEnabled {
			log.Println(UsedLogger.phosphorDbg("[🔍 DEBUG]") + " " + message)
		} else {
			log.Println("[DEBUG] " + message)
		}
	}
}

func Info(message string) {
	if UsedLogger.logLevel <= LogLevelINFO {
		if UsedLogger.ansiEnabled {
			log.Println(UsedLogger.phosphorInfo("[🛈 INFO]") + " " + message)
		} else {
			log.Println("[INFO] " + message)
		}
	}
}

func Output(message string) {
	if UsedLogger.logLevel <= LogLevelINFO {
		if UsedLogger.ansiEnabled {
			log.Println(UsedLogger.phosphorOutput("[🚂 OUTPUT]") + " " + message)
		} else {
			log.Println("[OUTPUT] " + message)
		}
	}
}

func Ok(message string) {
	if UsedLogger.logLevel <= LogLevelWARN {
		if UsedLogger.ansiEnabled {
			log.Println(UsedLogger.phosphorOk("[✅ OK]") + " " + message)
		} else {
			log.Println("[OK] " + message)
		}
	}
}

func Warn(message string) {
	if UsedLogger.logLevel <= LogLevelWARN {
		if UsedLogger.ansiEnabled {
			log.Println(UsedLogger.phosphorWarn("[⚠️ WARN]") + " " + message)
		} else {
			log.Println("[WARN] " + message)
		}
	}
}

func Err(message string) {
	if UsedLogger.logLevel <= LogLevelERROR {
		if UsedLogger.ansiEnabled {
			log.Println(UsedLogger.phosphorErr("[🔥 ERROR]") + " " + message)
		} else {
			log.Println("[ERROR] " + message)
		}
	}
}

func Critical(message string) {
	if UsedLogger.ansiEnabled {
		log.Fatalln(UsedLogger.phosphorCrit("[⚡ FATAL]") + " " + message)
	} else {
		log.Fatalln("[FATAL] " + message)
	}
}
