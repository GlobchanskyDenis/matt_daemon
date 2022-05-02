package main

import (
	"matt-daemon/pkg/netSocket"
	"os"
)

func main() {
	client := netSocket.NewTcpClient("localhost:8080")
	if err := client.Dial(); err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	client.SetReader(os.Stdin)
	client.SetWriter(os.Stdout)

	go func(client netSocket.Client){
		for {
			if err := client.ReadToPipe(); err != nil {
				println(err.Error())
				break
			}
		}
	}(client)

	for {
		if err := client.WriteToPipe(); err != nil {
			println(err.Error())
			break
		}
	}
}