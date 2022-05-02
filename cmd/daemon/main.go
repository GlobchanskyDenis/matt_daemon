package main

import (
	"matt-daemon/pkg/utils/file_logger"
	"matt-daemon/pkg/netSocket"
	"matt-daemon/pkg/utils/locker"
	"matt-daemon/pkg/constants"
	"os/signal"
	"syscall"
	"os"
)

func main() {
	if err := initializeConfigs("config/conf.json"); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	file_logger.GLogger.LogInfo("Стартую")
	file_logger.GLogger.LogInfo("Считан конфигурационный файл")

	if err := locker.Lock(constants.LockFilePath); err != nil {
		file_logger.GLogger.LogError(err, "Не могу запретить другим приложениям занимать порт")
		os.Exit(1)
	}
	defer func() {
		locker.Unlock()
		file_logger.GLogger.LogInfo("Удалил Lock файл")
	}()
	file_logger.GLogger.LogInfo("Создал Lock файл")

	server, err := newSocketServerByConfig(gConf)
	if err != nil {
		file_logger.GLogger.LogError(err, "Ошибка при открытии сокета")
		os.Exit(1)
	}
	file_logger.GLogger.LogInfo("Создал сокет сервер")

	exit := make(chan os.Signal, 1)

	go handleSocket(exit, server)

	waitForGracefullShutdown(exit, server)
}

func waitForGracefullShutdown(exit chan os.Signal, server netSocket.Server) {
	/*	Отлавливаю системный вызов останова программы. Это блокирующая операция  */
	signal.Notify(exit,
		syscall.SIGTERM, /*  Согласно всякой документации именно он должен останавливать прогу, но на деле его мы не находим. Оставил его просто на всякий случай  */
		syscall.SIGINT,  /*  Останавливает прогу когда она запущена из терминала и останавливается через CTRL+C  */
		syscall.SIGQUIT, /*  Останавливает демона systemd  */
	)
	<-exit

	if err := server.Close(); err != nil {
		file_logger.GLogger.LogError(err, "Ошибка во время закрытия сокета")
	}	
}