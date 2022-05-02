package netSocket

import (
	"net"
)

type tcpServer struct {
	address string
	conn    net.Conn
	buf     []byte
}

var _ Server = (*tcpServer)(nil)

func NewTcpServer(address string) Server {
	return &tcpServer{
		address: address,
		buf: make([]byte, 32),
	}
}

func (entity *tcpServer) Dial() error {
	listener, err := net.Listen("tcp", entity.address)
	if err != nil {
		return err
	}
	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	entity.conn = conn
	return nil
}

func (entity *tcpServer) Read() ([]byte, error) {
	if entity.conn != nil {
		readLen, err := entity.conn.Read(entity.buf)
		if err != nil {
			return nil, err
		}
		return entity.buf[:readLen], nil
	}
	return nil, nil
}

func (entity *tcpServer) Write(payload []byte) error {
	if entity.conn != nil {
		if _, err := entity.conn.Write(payload); err != nil {
			return nil
		}
	}
	return nil
}

func (entity *tcpServer) Close() error {
	if entity.conn != nil {
		if err := entity.conn.Close(); err != nil {
			return err
		}
		entity.conn = nil
	}
	return nil
}
