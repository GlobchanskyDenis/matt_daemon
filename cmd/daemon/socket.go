package main

import (
	"matt-daemon/pkg/utils/file_logger"
	"matt-daemon/pkg/netSocket"
	"matt-daemon/pkg/auth"
	"matt-daemon/pkg/dto"
	"strconv"
	"strings"
	"syscall"
	"os"
)

func newSocketServerByConfig(conf dto.SocketConfig) (netSocket.Server, error) {
	var server netSocket.Server
	if conf.IsTcpSocket == true {
		server = netSocket.NewTcpServer(conf.Ip + ":" + strconv.FormatUint(uint64(conf.Port), 10))
	} else {
		server = netSocket.NewUdpServer(conf.Ip, conf.Port)
	}

	if err := server.Dial(); err != nil {
		return nil, err
	}

	return server, nil
}

func handleSocket(exit chan os.Signal, server netSocket.Server) {
	if userIsAuthenticated(exit, server) {
		listenSocket(exit, server)
	}
}

func userIsAuthenticated(exit chan os.Signal, server netSocket.Server) bool {
	for {
		if err := server.Write([]byte("Введите логин\n")); err != nil {
			file_logger.GLogger.LogError(err, "Не смог авторизировать пользователя")
			exit<-syscall.SIGINT
			return false
		}
		loginRaw, err := server.Read()
		if err != nil {
			file_logger.GLogger.LogError(err, "Не смог авторизировать пользователя")
			exit<-syscall.SIGINT
			return false
		}
		login := parse(loginRaw)

		if err := server.Write([]byte("Введите пароль\n")); err != nil {
			file_logger.GLogger.LogError(err, "Не смог авторизировать пользователя")
			exit<-syscall.SIGINT
			return false
		}
		passwordRaw, err := server.Read()
		if err != nil {
			file_logger.GLogger.LogError(err, "Не смог авторизировать пользователя")
			exit<-syscall.SIGINT
			return false
		}
		password := parse(passwordRaw)

		if auth.IsExist(login, password) == true {
			file_logger.GLogger.LogInfo("user %s is authenticated", login)
			return true
		}
	}
}

func listenSocket(exit chan os.Signal, server netSocket.Server) {
	for {
		request, err := server.Read()
		if err != nil {
			file_logger.GLogger.LogError(err, "Ошибка чтения из сокета")
			exit<-syscall.SIGINT
			return
		}
		parsedRequest := parse(request)
		if isExitCommand(parsedRequest) {
			exit<-syscall.SIGINT
			return
		}

		if parsedRequest != "" {
			file_logger.GLogger.Log("user input: %s", parsedRequest)
		}
	}
}

func parse(request []byte) string {
	return strings.Trim(string(request), "	 \n")
}

func isExitCommand(src string) bool {
	if src == "exit" || src == "Exit" || src == "EXIT" {
		return true
	}
	return false
}