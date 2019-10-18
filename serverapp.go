package main

import (
	"fmt"
	"gitlab.dove.im/wx/cc_server_common/iface"
	"gitlab.dove.im/wx/cc_server_common/server"
)

func main() {

	//readChan := make(chan []byte , 100)

	s := server.NewServer("192.168.30.115", "20002")

	err := s.StartServer(func(conn iface.IConnection) {

		fmt.Println("New Conn  ConnID is ", conn.GetConnID())

		go func() {

			for {
				readData, err := conn.ReadData()

				if err != nil {
					fmt.Println("Read Err is ", err)
					return
					//continue
				}

				fmt.Printf("ReadData To Conn ID is %d \n", conn.GetConnID())

				//fmt.Println("---------------Server ReadData 100 is ", readData[:100])
				//fmt.Println("---------------Server ReadData -100 is ", readData[len(readData)-101:])
				//time.Sleep(1 * time.Second)
				// 立即回复
				conn.WriteData(readData)
			}

		}()

		//conn.WriteData([]byte("qweqwe"))
		fmt.Println("New Conn  Conn Remote is ", conn.GetRemoteAdd())
	}, nil)

	if err != nil {
		fmt.Println(err)
	}

}
