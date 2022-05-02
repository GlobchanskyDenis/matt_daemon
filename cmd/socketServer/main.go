package main

import (
	"matt-daemon/pkg/netSocket"
	"os"
)

func main() {
	server := netSocket.NewTcpServer("localhost:8080")
	if err := server.Dial(); err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	for {
		if err := server.Write([]byte("Hello, say something\n")); err != nil {
			println(err.Error())
			os.Exit(-1)
		}
		message, err := server.Read()
		if err != nil {
			println(err.Error())
			os.Exit(-1)
		}
		response := append([]byte("You said "), message...)
		if err := server.Write(response); err != nil {
			println(err.Error())
			os.Exit(-1)
		}
	}
}