package tcpserver

import (
	"fmt"
	"gitlab.dove.im/wx/cc_server_common/connMgr"
	"gitlab.dove.im/wx/cc_server_common/iface"
	"net"
)

type TcpServer struct {
	Ip          string
	Port        string
	Name        string
	Domain      string
	ConnManager *connMgr.ConnMgr

	//exitSigChan chan os.Signal

	isClose bool
}

func NewTcpServer(ip, port string) iface.IServer {
	return &TcpServer{
		Ip:          ip,
		Port:        port,
		Name:        "TcpServer",
		ConnManager: connMgr.NewConnMgr(),
		//exitSigChan: make(chan os.Signal),
		isClose: false,
	}
}

func (this *TcpServer) StartServer(newConnHandle iface.NewConnHandle, newMsgHandle iface.NewMessageHandle) error {

	curAddrs := fmt.Sprintf("%s:%s", this.Ip, this.Port)

	fmt.Println("Server Address : ", curAddrs)

	tcpAddr, err := net.ResolveTCPAddr("tcp", curAddrs)
	if err != nil {
		fmt.Println("ResolveTCPAddr Tcp Start Error ", err)
		return err
	}

	tcpLis, err := net.ListenTCP("tcp4", tcpAddr)

	if err != nil {
		fmt.Println("Server Tcp Start Error ", err)
		return err
	}

	var n int64
	for {

		if this.isClose {
			return nil
		}

		conn, err := tcpLis.AcceptTCP()

		if err != nil {
			fmt.Println("accept tcp err : ", err)
			continue
		}
		n++
		//NewConnection
		tcpConn := connMgr.NewConnection(this, conn, n, nil)
		this.ConnManager.AddConn(n, tcpConn)

		tcpConn.Start()

		newConnHandle(tcpConn)
	}

}

func (this *TcpServer) StopServer() {
	this.isClose = true

	this.ConnManager.StopAllConn()

}

func (this *TcpServer) GetConnMgr() iface.IConnMgr {
	return this.ConnManager
}

func (this *TcpServer) GetServerPort() string {
	return this.Port
}
