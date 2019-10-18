package server

import (
	"fmt"
	"github.com/xtaci/kcp-go"
	"gitlab.dove.im/wx/cc_server_common/connMgr"
	"gitlab.dove.im/wx/cc_server_common/iface"
)

//
const (
	ReadDeadTime  = 6
	WriteDeadTime = 6
)

type Server struct {
	Ip          string
	Port        string
	Name        string
	Domain      string
	ConnManager *connMgr.ConnMgr

	//exitSigChan chan os.Signal

	isClose bool
}

func NewServer(ip, port string) iface.IServer {
	return &Server{
		Ip:          ip,
		Port:        port,
		Name:        "KcpServer",
		ConnManager: connMgr.NewConnMgr(),
		//exitSigChan: make(chan os.Signal),
		isClose: false,
	}
}

func NewServerByDomain(domain, port string) *Server {
	return &Server{
		Domain:      domain,
		Port:        port,
		Name:        "KcpServer",
		ConnManager: connMgr.NewConnMgr(),
	}
}

// NewConnHandle
func (this *Server) StartServer(newConnHandle iface.NewConnHandle, newMsgHandle iface.NewMessageHandle) error {

	curAdd := ""
	if this.Ip == "" {
		curAdd = fmt.Sprintf("%s:%s", this.Domain, this.Port)
	} else {
		curAdd = fmt.Sprintf("%s:%s", this.Ip, this.Port)
	}
	//curAdd := fmt.Sprintf("%s:%s" , this.Ip , this.Port)
	fmt.Println("Cur Address is ", curAdd)

	// 监听 信号
	//signal.Notify(this.exitSigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	lis, err := kcp.ListenWithOptions(curAdd, nil, 10, 3)

	//lis, err := kcp.Listen(curAdd)

	if err != nil {
		fmt.Println("ServerBegin", err)
		return nil
	}
	var n int64
	for {

		if this.isClose {
			return nil
		}

		curConn, err := lis.Accept()

		if udp, ok := curConn.(*kcp.UDPSession); ok {
			udp.SetWriteBuffer(64 * 1024)
			udp.SetNoDelay(1, 10, 2, 1)
		}

		if err != nil {
			fmt.Println("Accept err ", err)
			continue
		}
		n++
		con := connMgr.NewConnection(this, curConn, n, newMsgHandle)
		this.ConnManager.AddConn(n, con)
		con.Start()

		newConnHandle(con)

	}

}
func (this *Server) StopServer() {

	this.isClose = true

	this.ConnManager.StopAllConn()
}

func (this *Server) GetConnMgr() iface.IConnMgr {
	return this.ConnManager
}

func (this *Server) GetServerPort() string {
	return this.Port
}

//func (this *Server) WriteData(connID int64, data []byte) error {
//
//	conn := this.ConnManager.GetConn(connID)
//
//	if conn == nil {
//		return errors.New("Get Connection Fail")
//	}
//
//	conn.WriteData(data)
//	return nil
//}
