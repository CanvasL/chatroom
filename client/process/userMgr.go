package process

import (
	"fmt"
	"go_socket/client/message"
	"go_socket/client/model"
)

// 客户端要维护的map
var onlineUsers map[int]*model.User = make(map[int]*model.User, 10)

// showOnlineUsers 在客户端显示当前在线的用户
func showOnlineUsers() {
	// 遍历onlineUsers即可
	for id, user := range onlineUsers {
		fmt.Println("用户id: ", id, "\t状态: ", user.UserStatus)
	}
}

// updateUserStatus 处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	// 适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		// 原来没有
		user = &model.User{
			UserId: notifyUserStatusMes.UserId,
			//UserStatus: notifyUserStatusMes.Status,
		}
	}
	// 如果原来就有，那么id不需要变化，只需要改变状态即可
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	showOnlineUsers()
}
