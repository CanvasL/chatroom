package process

import (
	"encoding/json"
	"fmt"
	"go_socket/client/message"
)

func printGroupMes(mes *message.Message) {
	// 这里来的类型一定是message.SmsMesType
	// 显示即可
	// 1.反序列化
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &smsMes) failed, err:", err)
		return
	}

	// 2.显示信息
	info := fmt.Sprintf("%s(id:%d):\t%s", smsMes.UserName, smsMes.UserId, smsMes.Content)
	fmt.Println(info)
}
