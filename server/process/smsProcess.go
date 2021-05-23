package process

import (
	"encoding/json"
	"fmt"
	"go_socket/server/message"
	"go_socket/server/utils"
	"net"
)

type SmsProcess struct {
	// 无字段，纯用来调用函数方法
}

// SendGroupMes 转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	// 遍历服务器端的onlineUsers，将消息转发

	// 取出mes的内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &smsMes) failed, err:", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) failed, err:", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		// 过滤掉自己
		if id == smsMes.UserId {
			continue
		}
		this.SendToEveryOnlineUser(data, up.Conn)
	}
}

// SendToEveryOnlineUser 把消息发送给每个在线用户
func (this *SmsProcess) SendToEveryOnlineUser(data []byte, conn net.Conn) {

	// 创建一个Transfer实例，发送给data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败！err:", err)
	}
}
