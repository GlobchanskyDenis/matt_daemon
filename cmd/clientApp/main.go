package main

import (
	"matt-daemon/pkg/window"
)

func main() {
	win := window.New(isConnectServerCorrect, isAuthCorrect, sendMessage)
	defer socketClose()
	win.ShowAndRun()
}

