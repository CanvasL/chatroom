package process

import (
	"encoding/json"
	"fmt"
	"go_socket/client/message"
	"go_socket/client/model"
	"go_socket/client/utils"
	"net"
	"os"
)

// ShowMenu 显示登录成功后的界面
func ShowMenu() {
	//fmt.Println("=>>恭喜xxx登录成功。")
	fmt.Printf("========================%s=========================\n", model.CurrentUser.UserName)
	fmt.Printf("  用户id: %d\n", model.CurrentUser.UserId)
	fmt.Printf("  状  态: %s\n", model.CurrentUser.UserStatus)
	fmt.Println("\t\t1.显示用户在线列表")
	fmt.Println("\t\t2.发送消息")
	fmt.Println("\t\t3.信息列表")
	fmt.Println("\t\t4.退出系统")
	fmt.Println("====二级目录========")

	fmt.Println("Sys>>请选择(1~4):")
	fmt.Printf(">")

	var key int
	var content string

	// 因为我们总会使用到SmsProcess实例，因此我们将其定义在switch外部
	fmt.Scanf("%d\n", &key)
	smsProcess := &SmsProcess{}

	switch key {
	case 1:
		fmt.Println("=>>显示在线用户列表")
		fmt.Println("========================在线用户=========================")
		showOnlineUsers()
		fmt.Println("====三级目录========")

	case 2:
		fmt.Println("=>>发送消息")
		fmt.Println("========================群聊=========================")
		fmt.Println("〇回车键发送，输入exit退出群聊。")
		for content != "exit" {
			fmt.Printf(">")
			fmt.Scanf("%s\n", &content)
			smsProcess.SendGroupMes(content)
		}
		fmt.Println("====三级目录========")

	case 3:
		fmt.Println("=>>信息列表")

	case 4:
		fmt.Println("=>>你选择退出了系统。")
		os.Exit(0)

	default:
		fmt.Println(">>你输入的选项不正确。")
	}
}

// processServerMes 保持和server的通讯，如有信息则显示在客户端
func processServerMes(conn net.Conn) {
	// 创建一个Transfer实例，不停地读取服务端发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		//fmt.Println("客户端%s正在等待读取服务器发送的消息...")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("读取服务器端出错, err:", err)
			return
		}
		// 如果读取到消息，又是下一步处理逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 有人上线
			fmt.Printf("\nSvr>>有新的好友上线了!\n")
			// 1.取出 NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			// 2.把这个用户的信息、状态保存到客户map中
			updateUserStatus(&notifyUserStatusMes)

		case message.SmsMesType:
			// 有人群发消息
			fmt.Printf("\nSvr>>有新的群聊消息!\n")
			printGroupMes(&mes)

		default:
			fmt.Printf("\nSvr>>错误，服务器返回了未知消息。\n")
		}
	}
}
