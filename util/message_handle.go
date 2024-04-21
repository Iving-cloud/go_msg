package util

import (
	"fmt"
	"go_msg/proto_msg"
	"net"

	"google.golang.org/protobuf/proto"
)

func HelloSend(worker_id int32, conn net.Conn) {
	buf := make([]byte, 1024)

	hello_msg := &proto_msg.HeartHello{
		Type:proto_msg.MonitorType_HELLO,
		WorkerId: worker_id,
	}
	buf,_ = proto.Marshal(hello_msg)
	conn.Write([]byte(buf))
}


func Decode_msg(buf []byte) proto_msg.MonitorType{
	msg := &proto_msg.BaseMsg{}
	proto.Unmarshal(buf, msg)
	fmt.Println("MSG type:", msg.GetType())
	return msg.GetType()
}