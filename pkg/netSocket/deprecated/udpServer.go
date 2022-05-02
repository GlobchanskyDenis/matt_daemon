package main

import (
   "fmt"
   "net"
)

func main() {
   listener, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("localhost"), Port: 8080 }) // открываем слушающий UDP-сокет
   for {
      handleClient(listener) // обрабатываем запрос клиента
   }
}

func handleClient(conn *net.UDPConn) {
   buf := make([]byte, 128) // буфер для чтения клиентских данных
   
   readLen, addr, err := conn.ReadFromUDP(buf) // читаем из сокета
   if err != nil {
      fmt.Println(err)
      return
   }

   conn.WriteToUDP(append([]byte("Hello, you said: "), buf[:readLen]...), addr) // пишем в сокет
}