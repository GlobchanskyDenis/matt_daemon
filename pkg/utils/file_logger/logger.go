package file_logger

import (
	// "10.10.11.220/ursgis/cdocs_epgu_sender_receiver.git/pkg/utils"
	"fmt"
	"os"
	"sync"
	"time"
)

var GLogger Logger

/*	мьютекс необходим так как под капотом логгер может менять файл и необходимо обеспечить потокобезопасность  */
type logger struct {
	filePath        string
	fileName        string
	permissions     string
	currentFileDate time.Time
	fileOS          *os.File
	mu              *sync.Mutex
}

type Logger interface {
	LogError(err error, format string, args ...interface{})
	LogWarning(err error, format string, args ...interface{})
	LogInfo(format string, args ...interface{})
	Log(format string, args ...interface{})
	LogDebug(format string, args ...interface{})
}

func NewLogger() error {
	if err := checkConfig(); err != nil {
		return err
	}
	conf := GetConfig()

	newLogger := &logger{
		filePath:    conf.LogFolder,
		permissions: conf.Permissions,
		mu:          &sync.Mutex{},
	}

	if err := newLogger.setNewLogFile(); err != nil {
		return err
	}

	GLogger = newLogger

	return nil
}

func (l *logger) LogError(err error, format string, args ...interface{}) {
	if newErr := l.changeLogFileIfItNeeded(); newErr != nil {
		fmt.Printf("При смене лог файла %s %#v %s\n", newErr, err, fmt.Sprintf(format, args...))
	} else {
		message := prepareToLogThis("ERROR", err, format, args...)
		/*	Повторная попытка нужна для случая когда дата сменилась и какой-то воркер
		**	в другом потоке начал менять дескриптор файла для логгирования. Считаю секунды
		**	достаточно. Не блокировать же мьютексами запись в файл :)  */
		if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
			time.Sleep(1 * time.Second)
			if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
				print(err.Error() + " => " + message)
			}
		}
	}
}

func (l *logger) LogWarning(err error, format string, args ...interface{}) {
	if newErr := l.changeLogFileIfItNeeded(); newErr != nil {
		fmt.Printf("При смене лог файла %s %#v %s\n", newErr, err, fmt.Sprintf(format, args...))
	} else {
		message := prepareToLogThis("WARNING", err, format, args...)
		/*	Повторная попытка нужна для случая когда дата сменилась и какой-то воркер
		**	в другом потоке начал менять дескриптор файла для логгирования. Считаю секунды
		**	достаточно. Не блокировать же мьютексами запись в файл :)  */
		if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
			time.Sleep(1 * time.Second)
			if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
				print(err.Error() + " => " + message)
			}
		}
	}
}

func (l *logger) LogInfo(format string, args ...interface{}) {
	if newErr := l.changeLogFileIfItNeeded(); newErr != nil {
		fmt.Printf("При смене лог файла %s %s\n", newErr, fmt.Sprintf(format, args...))
	} else {
		message := prepareToLogThis("INFO", nil, format, args...)
		/*	Повторная попытка нужна для случая когда дата сменилась и какой-то воркер
		**	в другом потоке начал менять дескриптор файла для логгирования. Считаю секунды
		**	достаточно. Не блокировать же мьютексами запись в файл :)  */
		if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
			time.Sleep(1 * time.Second)
			if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
				print(err.Error() + " => " + message)
			}
		}
	}
}

func (l *logger) Log(format string, args ...interface{}) {
	if newErr := l.changeLogFileIfItNeeded(); newErr != nil {
		fmt.Printf("При смене лог файла %s %s\n", newErr, fmt.Sprintf(format, args...))
	} else {
		message := prepareToLogThis("LOG", nil, format, args...)
		/*	Повторная попытка нужна для случая когда дата сменилась и какой-то воркер
		**	в другом потоке начал менять дескриптор файла для логгирования. Считаю секунды
		**	достаточно. Не блокировать же мьютексами запись в файл :)  */
		if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
			time.Sleep(1 * time.Second)
			if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
				print(err.Error() + " => " + message)
			}
		}
	}
}

func (l *logger) LogDebug(format string, args ...interface{}) {
	if newErr := l.changeLogFileIfItNeeded(); newErr != nil {
		fmt.Printf("При смене лог файла %s %s\n", newErr, fmt.Sprintf(format, args...))
	} else {
		message := prepareToLogThis("DEBUG", nil, format, args...)
		/*	Повторная попытка нужна для случая когда дата сменилась и какой-то воркер
		**	в другом потоке начал менять дескриптор файла для логгирования. Считаю секунды
		**	достаточно. Не блокировать же мьютексами запись в файл :)  */
		if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
			time.Sleep(1 * time.Second)
			if _, err := fmt.Fprintf(l.fileOS, "%s", message); err != nil {
				print(err.Error() + " => " + message)
			}
		}
	}
}

func prepareToLogThis(level string, err error, format string, args ...interface{}) string {
	prepared := fmt.Sprintf(format, args...)
	if err != nil {
		prepared += " " + err.Error()
	}
	var name string = "Matt_daemon"
	if gConf != nil {
		name = gConf.DaemonName
	}
	now := time.Now()
	return fmt.Sprintf("[%02d/%02d/%d-%02d:%02d:%02d] [ %s ] - %s: %s.\n", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second(), level, name, prepared)
}

/*	Меняет файл в который записывается логгирование в случае если уже сменилась дата
**	Использует мьютекс, поэтому выполняется горутинобезопасно  */
func (l *logger) changeLogFileIfItNeeded() error {
	if isSameDate(l.currentFileDate, time.Now()) == false {
		l.mu.Lock()
		defer l.mu.Unlock()
		if err := l.setNewLogFile(); err != nil {
			return err
		}
	}
	return nil
}

func isSameDate(oldDate, now time.Time) bool {
	if oldDate.Year() == now.Year() && oldDate.Month() == now.Month() && oldDate.Day() == now.Day() {
		return true
	}
	return false
}

/*	Данную функцию в многопоточном режиме нужно запускать в   */
func (l *logger) setNewLogFile() error {
	/*	Если ранее уже был открыт файл - закрываю его  */
	if l.fileOS != nil {
		if err := l.fileOS.Close(); err != nil {
			return err
		}
	}

	if err := checkConfig(); err != nil {
		return err
	}
	conf := GetConfig()

	/*	Открываю новый файл  */
	l.currentFileDate = time.Now()
	l.fileName = fmt.Sprintf("%s_%d-%d-%d.log", conf.DaemonName, l.currentFileDate.Year(), l.currentFileDate.Month(), l.currentFileDate.Day())
	osFile, err := openOrCreateNewFile(l.filePath, l.fileName, l.permissions)
	if err != nil {
		return err
	}
	l.fileOS = osFile
	return nil
}
