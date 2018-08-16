package tcp

import (
	"log"
	"posam/util/concurrentmap"
	"time"
)

type TCPClienter interface {
	Send([]byte, []byte) ([]byte, error)
}

type TCPClient struct {
	Connectivitier
	ServerNetwork     string
	ServerAddress     string
	ServerTimeout     time.Duration
	ServerConcurrency bool
}

func NewTCPClient(network string, address string, second int, concurrency bool) *TCPClient {
	timeout := time.Duration(second) * time.Second
	return &TCPClient{
		Connectivitier:    NewConnectivity(network, address, timeout),
		ServerNetwork:     network,
		ServerAddress:     address,
		ServerTimeout:     timeout,
		ServerConcurrency: concurrency,
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	connectivityMap = concurrentmap.NewConcurrentMap()
}

type TCPRequest struct {
	Message  []byte
	Expected []byte
	Respc    chan TCPResponse
}

type TCPResponse struct {
	Response []byte
	Error    error
}

func (t *TCPClient) Send(message []byte, expected []byte) (resp []byte, err error) {
	req := TCPRequest{
		Message:  message,
		Expected: expected,
		Respc:    make(chan TCPResponse),
	}
	connectivity := t.Connectivitier.(*Connectivity)
	connectivity.RequestQueue.Push(req)
	log.Println("waiting for response...", message)
	tcpResponse := <-req.Respc
	return tcpResponse.Response, tcpResponse.Error
}
