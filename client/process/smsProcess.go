package process

import (
	"encoding/json"
	"fmt"
	"go_socket/client/message"
	"go_socket/client/model"
	"go_socket/client/utils"
)

type SmsProcess struct {
}

// 发送群聊的消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	// 1.创建一个mes
	var mes message.Message
	mes.Type = message.SmsMesType

	// 2.创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = model.CurrentUser.UserId
	smsMes.UserName = model.CurrentUser.UserName
	smsMes.UserStatus = model.CurrentUser.UserStatus

	// 3.序列化 smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("json.Marshal smsMes failed, err:", err)
		return
	}
	mes.Data = string(data)

	// 4.对 mes 进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal sms failed, err:", err)
		return
	}

	// 5.将data发送给服务器
	tf := &utils.Transfer{
		Conn: model.CurrentUser.Conn,
	}

	// 6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("tf.WritePkg data failed, err:", err)
		return
	}
	return
}
