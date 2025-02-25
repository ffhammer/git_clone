package logging

import (
	"errors"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/settings"
	"git_clone/gvc/utils"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type LoggerType struct {
	file  *os.File
	level settings.LogLevel
	log   *log.Logger
}

var (
	Logger  LoggerType
	once    sync.Once
	initErr error
)

func initLogger() error {
	once.Do(func() {

		logPath := filepath.Join(utils.RepoDir, config.LOG_PATH)

		setting, err := settings.LoadSettings()
		if err != nil {
			initErr = fmt.Errorf("error initializing logger: %w", err)
			return
		}

		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			initErr = fmt.Errorf("error opening log file: %w", err)
			return
		}

		logger := log.New(file, "", log.LstdFlags|log.Lmicroseconds)

		Logger = LoggerType{
			file:  file,
			level: setting.LogLevel,
			log:   logger,
		}
		initErr = nil
	})

	return initErr
}

func shouldLog(current, msg settings.LogLevel) bool {
	levels := map[settings.LogLevel]int{
		settings.DEBUGGING: 1,
		settings.INFO:      2,
		settings.WARNING:   3,
		settings.ERROR:     4,
	}

	return levels[msg] >= levels[current]
}

func log_message(message string, level settings.LogLevel) {
	if err := initLogger(); err != nil {
		fmt.Printf("can't log because of logger error: %v\n", err)
		return
	}

	if shouldLog(Logger.level, level) {
		Logger.log.Printf("[%s] %s", level, message)
	}
}

func ErrorF(format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	log_message(err.Error(), settings.ERROR)
	return err
}

func Error(err error) error {
	log_message(err.Error(), settings.ERROR)
	return err
}

func NewError(desc string) error {
	err := errors.New(desc)
	log_message(desc, settings.ERROR)
	return err
}

func Debug(message string) {
	log_message(message, settings.DEBUGGING)
}

func DebugF(format string, a ...any) {
	log_message(fmt.Sprintf(format, a...), settings.DEBUGGING)
}

func Info(message string) {
	log_message(message, settings.INFO)
}
func InfoF(format string, a ...any) {
	log_message(fmt.Sprintf(format, a...), settings.INFO)
}

func Warn(message string) {
	log_message(message, settings.WARNING)
}

func WarnF(format string, a ...any) {
	log_message(fmt.Sprintf(format, a...), settings.WARNING)
}
