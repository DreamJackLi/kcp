package iface

type IClient interface {
	StartClient(ip, port string) error
	StopClient() error
	ReadData() ([]byte, error)
	WriteData([]byte) error
}
