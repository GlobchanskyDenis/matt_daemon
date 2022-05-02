package main

import (
	"matt-daemon/pkg/constants"
	"matt-daemon/pkg/auth"
	"matt-daemon/pkg/dto"
	"matt-daemon/pkg/utils/file_logger"
	"matt-daemon/pkg/utils/u_conf"
)

var gConf dto.SocketConfig

func initializeConfigs(configFileName string) error {
	/*	Read config file  */
	print("Считываю конфигурационный файл\t- ")
	if err := u_conf.SetConfigFile(configFileName); err != nil {
		println(constants.RED + "ошибка" + constants.NO_COLOR)
		return err
	}
	println(constants.GREEN + "успешно" + constants.NO_COLOR)

	/*	socket  */
	print("настраиваю конфигурационник работы с сокетами\t- ")
	if err := u_conf.ParsePackageConfig(&gConf, "Socket"); err != nil {
		println(constants.RED + "ошибка" + constants.NO_COLOR)
		return err
	}
	println(constants.GREEN + "успешно" + constants.NO_COLOR)

	/*	auth  */
	print("настраиваю конфигурационник работы с авторизацией\t- ")
	authConf := auth.GetConfig()
	if err := u_conf.ParsePackageConfig(authConf, "Authenticate"); err != nil {
		println(constants.RED + "ошибка" + constants.NO_COLOR)
		return err
	}
	println(constants.GREEN + "успешно" + constants.NO_COLOR)

	/*	file_logger  */
	print("настраиваю пакет file_logger\t- ")
	loggerConf := file_logger.GetConfig()
	if err := u_conf.ParsePackageConfig(loggerConf, "Logger"); err != nil {
		println(constants.RED + "ошибка" + constants.NO_COLOR)
		return err
	}
	if err := file_logger.NewLogger(); err != nil {
		println(constants.RED + "ошибка" + constants.NO_COLOR)
		return err
	}
	println(constants.GREEN + "успешно" + constants.NO_COLOR)

	return nil
}
