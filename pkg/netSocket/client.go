package netSocket

import (
	"io"
	"net"
)

var _ Client = (*client)(nil)

type Client interface {
	Dial() error
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
	SetReader(reader io.Reader)
	SetWriter(writer io.Writer)
	ReadToPipe() error
	WriteToPipe() error
}

type client struct {
	network string
	address string
	conn    net.Conn
	buf     []byte
	reader  io.Reader
	writer  io.Writer
}

func NewTcpClient(address string) Client {
	return &client{
		network: "tcp",
		address: address,
		buf: make([]byte, 32),
	}
}

func NewUdpClient(address string) Client {
	return &client{
		network: "udp",
		address: address,
		buf: make([]byte, 32),
	}
}

func (entity *client) Dial() error {
	conn, err := net.Dial(entity.network, entity.address)
	if err != nil {
		return err
	}
	entity.conn = conn
	return nil
}

func (entity *client) Read() ([]byte, error) {
	if entity.conn != nil {
		readLen, err := entity.conn.Read(entity.buf)
		if err != nil {
			return nil, err
		}
		return entity.buf[:readLen], nil
	}
	return nil, nil
}

func (entity *client) Write(payload []byte) error {
	if entity.conn != nil {
		if _, err := entity.conn.Write(payload); err != nil {
			return err
		}
	}
	return nil
}

func (entity *client) Close() error {
	if entity.conn != nil {
		if err := entity.conn.Close(); err != nil {
			return err
		}
		entity.conn = nil
	}
	return nil
}

func (entity *client) SetReader(reader io.Reader) {
	entity.reader = reader
}

func (entity *client) SetWriter(writer io.Writer) {
	entity.writer = writer
}

func (entity *client) ReadToPipe() error {
	if entity.reader != nil && entity.conn != nil {
		if _, err := io.Copy(entity.conn, entity.reader); err != nil {
			return err
		}
	}
	return nil
}

func (entity *client) WriteToPipe() error {
	if entity.writer != nil && entity.conn != nil {
		if _, err := io.Copy(entity.writer, entity.conn); err != nil {
			return err
		}
	}
	return nil
}
