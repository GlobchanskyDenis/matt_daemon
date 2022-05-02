package main

import (
	"matt-daemon/pkg/utils/file_logger"
	"matt-daemon/pkg/utils/process"
	"matt-daemon/pkg/utils/locker"
	"matt-daemon/pkg/constants"
	"os"
)

func main() {
	if err := initializeConfigs("config/conf.json"); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if locker.IsLocked(constants.LockFilePath) == true {
		file_logger.GLogger.LogInfo("Не имею права запускаться так как Lock файл уже был создан")
		println("Не имею права запускаться так как Lock файл уже был создан")
		os.Exit(1)
	}

	/*	Создаем процесс  */
	cmd, _, _, err := process.New("./daemon_bin", "./", nil, nil)
	if err != nil {
		file_logger.GLogger.LogError(err, "Ошибка при запуске процесса")
		println(err.Error())
		os.Exit(1)
	}

	file_logger.GLogger.LogInfo("Поднял процесс на порту %d", int(cmd.Process.Pid))
}