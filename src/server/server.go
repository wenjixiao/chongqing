package main

import (
    "os"
	"fmt"
	"io"
	"net"
	"bytes"
	"encoding/binary"
	)
//=========================================================
const HeaderSize = 4

type MsgHandler interface {
    Connected(conn net.Conn)
	Process(conn net.Conn,msgBytes []byte)
	Closed(conn net.Conn)
	Exception(conn net.Conn,err error)
}

//=========================================================
func AddHeader(msgBytes []byte) []byte {
	head := make([]byte, HeaderSize)
	binary.LittleEndian.PutUint32(head, uint32(len(msgBytes)))
	return append(head, msgBytes...)
}

//=========================================================
/*
	1,keep reading data from socket 
	2,parse data to msg
	3,give msg to MsgHandler
*/
func ReadMsg(conn net.Conn,handler MsgHandler) {
    
    handler.Connected(conn)
    
	defer conn.Close()

	const MSG_BUF_LEN = 1024 * 9
	const READ_BUF_LEN = 1024       //1KB

	fmt.Printf("Client: %s\n", conn.RemoteAddr())

	msgBuf := bytes.NewBuffer(make([]byte, 0, MSG_BUF_LEN))
	readBuf := make([]byte, READ_BUF_LEN)

	head := uint32(0)
	bodyLen := 0 //bodyLen is a flag,when readed head,but body'len is not enougth

	for {
		n, err := conn.Read(readBuf)
		if err != nil {
			if err == io.EOF {
				handler.Closed(conn)
			} else {
				handler.Exception(conn,err)
			}
			//close the connection
			break
		}
		_, err = msgBuf.Write(readBuf[:n])

		if err != nil {
		    fmt.Println("Buffer write error: ", err)
		    os.Exit(1)
		}

		for {
			//read the msg head
			if bodyLen == 0 && msgBuf.Len() >= HeaderSize {
				err := binary.Read(msgBuf, binary.LittleEndian, &head)
				if err != nil {
					fmt.Println("msg head Decode error: ", err)
				}
				bodyLen = int(head)

				if bodyLen > MSG_BUF_LEN {
					fmt.Println("msg body too long: ", bodyLen)
					os.Exit(1)
				}
			}
			//has head,now read body
			if bodyLen > 0 && msgBuf.Len() >= bodyLen {
				handler.Process(conn,msgBuf.Next(bodyLen))
				bodyLen = 0
			} else {
				//msgBuf.Len() < bodyLen ,one msg receiving is not complete
				//need to receive again
				break
			}
		} //for of msg buf
	} //for of conn read
}

//=========================================================
/*
	1,keep accepting connection
	2,make a goroutine process the connection with the msgHandler
*/
func ListenLoop() {
    //listen socket
	l, err := net.Listen("tcp", ":5678")
	defer l.Close()
	if err != nil {
	    fmt.Println("socket listen error: ",err)
	    os.Exit(1)
	}
	
	//start the endless cycle
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//here,we create new shadow for every routine
		go ReadMsg(conn,&Shadow{})
	}
}

//=========================================================
//services
var shadows_service ShadowsService
var games_service GamesService
var linebroken_service LinebrokenService
//=========================================================

func main(){
    //first,start all services
    shadows_service.start()
    games_service.start()
    linebroken_service.start()

	ListenLoop()
	//i don't want make a exe file,so when compile to soso,must be error!
	soso()
}