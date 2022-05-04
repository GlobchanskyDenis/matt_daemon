package main

import (
	"matt-daemon/pkg/netSocket"
	"strconv"
	"strings"
	"fmt"
)

var gClient netSocket.Client

func isConnectServerCorrect(isTcpConn bool, ip string, port uint) bool {
	fmt.Printf("is tcp connection %#v ip %s port %d\n", isTcpConn, ip, int(port))

	address := ip + ":" + strconv.FormatUint(uint64(port), 10)

	if isTcpConn == true {
		gClient = netSocket.NewTcpClient(address)
	} else {
		gClient = netSocket.NewUdpClient(address)
	}

	if err := gClient.Dial(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	/*	Нужно пропустить первую строчку приглашения сервера "введите логин" чтобы авторизоваться в цикле  */
	if _, err := gClient.Read(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	return true
}

func isAuthCorrect(login, password string) bool {
	fmt.Printf("login %s password %s\n", login, password)

	if err := gClient.Write([]byte(login)); err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	
	/*	Нужно пропустить первую строчку приглашения сервера "введите пароль" чтобы авторизоваться в цикле  */
	if _, err := gClient.Read(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}

	if err := gClient.Write([]byte(password)); err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}

	/*	Cчитываю строчку и анализирую по ней - выполнили ли мы авторизацию или нет. Если фраза
	**	"Авторизация успешна" то выполнили, если фраза - 
	**	"Введите логин" - то нет  */
	lineRaw, err := gClient.Read()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	fmt.Println(string(lineRaw))
	if strings.HasPrefix(string(lineRaw), "Авторизация усп") == false {
		return false
	}
	return true
}

func parse(request []byte) string {
	return strings.Trim(string(request), "	 \n")
}

func sendMessage(message string) {
	if err := gClient.Write([]byte(message)); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func socketClose() {
	if gClient != nil {
		if err := gClient.Close(); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		gClient = nil
	}
}
