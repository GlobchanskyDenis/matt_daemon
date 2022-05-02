package main

import (
   "fmt"
   "net"
)

func main() {
   listener, _ := net.Listen("tcp", "localhost:8080") // открываем слушающий сокет
   for {
      conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
      if err != nil {
         continue
      }
      go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
   }
}

func handleClient(conn net.Conn) {
   defer conn.Close() // закрываем сокет при выходе из функции

   buf := make([]byte, 32) // буфер для чтения клиентских данных
   for {
      conn.Write([]byte("Hello, what's your name?\n")) // пишем в сокет

      readLen, err := conn.Read(buf) // читаем из сокета
      if err != nil {
         fmt.Println(err)
         break
      }

      conn.Write(append([]byte("Goodbye, "), buf[:readLen]...)) // пишем в сокет
   }
}
