package netSocket

import (
	"net"
)

/*	ВАЖНО !!
**	lastAddress в составе структурки - из-за нее не предполагается работа в многопоточном режиме!! */

type udpServer struct {
	ip          string
	port        uint
	buf         []byte
	listener    *net.UDPConn
	lastAddress *net.UDPAddr
}

var _ Server = (*udpServer)(nil)

func NewUdpServer(ip string, port uint) Server {
	return &udpServer{
		ip: ip,
		port: port,
		buf: make([]byte, 128),
	}
}

func (entity *udpServer) Dial() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(entity.ip), Port: int(entity.port)})
	if err != nil {
		return err
	}
	entity.listener = listener
	return nil
}

func (entity *udpServer) Read() ([]byte, error) {
	if entity.listener != nil {
		readLen, addr, err := entity.listener.ReadFromUDP(entity.buf)
		if err != nil {
			return nil, err
		}
		entity.lastAddress = addr
		return entity.buf[:readLen], nil
	}
	return nil, nil
}

func (entity *udpServer) Write(payload []byte) error {
	if entity.listener != nil && entity.lastAddress != nil {
		if _, err := entity.listener.WriteToUDP(payload, entity.lastAddress); err != nil {
			return nil
		}
	}
	return nil
}

func (entity *udpServer) Close() error {
	if entity.listener != nil {
		if err := entity.listener.Close(); err != nil {
			return err
		}
		entity.listener = nil
	}
	return nil
}
