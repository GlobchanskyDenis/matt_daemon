package netSocket

type Server interface {
	Dial() error
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}
