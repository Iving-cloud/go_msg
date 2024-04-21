package main

import (
	"fmt"
	"go_msg/proto_msg"
	"go_msg/util"
	"go_msg/workers"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"os"
	"time"
)

var workerId int32

// 定义一个channel，用于接收任务完成或超时的信号
var done chan bool

func main() {
	//var wg sync.WaitGroup
	// 监听本地端口
	listener, err := net.Listen("tcp", ":9090")
	workerId = 1
	done = make(chan bool)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on :9090")
	go workers.PrintWorkerInfo()
	for {
		// 接受新的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		var worker = workers.WorkerInit(conn, 0)
		// 处理连接
		go handleRequest(worker)
	}

}

func RegisterHandle(w workers.Worker) {
	fmt.Println("PID:", w.Pid)
	msg := &proto_msg.RegisterResponse{
		Type:     proto_msg.MonitorType_RESPONSE,
		WorkerId: int32(w.WorkerId),
		Pid:      w.Pid,
		Role:     w.Role,
	}
	buf, _ := proto.Marshal(msg)
	w.Conn.Write([]byte(buf))
}

func periodicTask() {
	fmt.Println(time.Now().Second(), "Timeout!")
}

// 处理请求的函数
func handleRequest(worker *workers.Worker) {

	conn := worker.Conn
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Println("err", err)
		}
		//fmt.Printf("<<<<<== server get client message:\n")

		msgType := util.Decode_msg(buf)
		switch msgType {
		case proto_msg.MonitorType_HELLO:
			helloMsg := &proto_msg.HeartHello{}
			proto.Unmarshal(buf, helloMsg)
			util.HelloSend(0, conn)
			//conn.Write([]byte("Receive Hello!"))
			worker.ChTimer <- 1
			fmt.Printf("Receive worker-id %d Hello!\n", helloMsg.GetWorkerId())
		case proto_msg.MonitorType_REGISTER:
			request := &proto_msg.RegisterRequest{}
			proto.Unmarshal(buf, request)
			var job util.MyJob
			job.W = worker
			worker.Pid = request.GetPid()
			workers.AddWorker(*worker)
			go util.MyTimer(job)
			RegisterHandle(*worker)
			worker.ChTimer <- 0
		default:
			fmt.Println("Invalid Type!")
		}
	}
}
