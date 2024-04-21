package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"go_msg/proto_msg"
	"go_msg/util"

	"google.golang.org/protobuf/proto"
)

var worker_id int32

func ReceiveHandle(conn net.Conn) {
	//reader := bufio.NewReader(conn)
	fmt.Println("Start Receive!")
	for {
		// 接收来自服务器的响应
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		//buffer, err := Decode(reader)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		msg_type := util.Decode_msg(buffer)
		switch msg_type {
		case proto_msg.MonitorType_HELLO:
			fmt.Println("Receive HELLO resp!")
		case proto_msg.MonitorType_RESPONSE:
			response := &proto_msg.RegisterResponse{}
			proto.Unmarshal(buffer, response)
			fmt.Println("Role:", response.GetRole())
			worker_id = response.GetWorkerId()
		default:
			fmt.Println("Invalid Type!")
		}
	}
}

func RegisterToServer(conn net.Conn) {
	worker := &proto_msg.RegisterRequest{
		Type: proto_msg.MonitorType_REGISTER,
		Pid:  int32(os.Getpid()),
	}
	fmt.Println("pid:", worker.GetPid())
	fmt.Printf("%#v\n", worker)

	buf, _ := proto.Marshal(worker)
	conn.Write([]byte(buf))
}

func say(conn net.Conn) {
	for {
		time.Sleep(5 * time.Second)
		util.HelloSend(worker_id, conn)
	}
}

func main() {
	var wg sync.WaitGroup
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	RegisterToServer(conn)

	wg.Add(1)
	go ReceiveHandle(conn)
	wg.Add(1)
	go say(conn)
	wg.Wait()
	defer conn.Close()
}
