package main

import (
	"fmt"
	"go_socket/client/process"
	"os"
)

// 定义两个变量，一个表示用户id，一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {

	// 接收用户的选择
	var key int
	// 判断是否还继续显示菜单
	//var loop bool = true

	for {
		fmt.Println("==============欢迎登录多人聊天系统==============")
		fmt.Println("\t\t 1 登录聊天室")
		fmt.Println("\t\t 2 注册用户")
		fmt.Println("\t\t 3 退出系统")
		fmt.Println("====一级目录========")
		fmt.Println("Sys>>请选择(1~3):")

		fmt.Printf(">")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("=>>登录聊天室")
			fmt.Println("Sys>>请输入用户的id:")
			fmt.Printf(">")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("Sys>>请输入用户的密码:")
			fmt.Printf(">")
			fmt.Scanf("%s\n", &userPwd)
			// 完成登录
			// 创建UserProcess的实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)

		case 2:
			fmt.Println("=>>注册用户")
			fmt.Println("Sys>>请输入用户id:")
			fmt.Printf(">")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("Sys>>请输入用户密码:")
			fmt.Printf(">")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("Sys>>请输入用户昵称:")
			fmt.Printf(">")
			fmt.Scanf("%s\n", &userName)
			// 完成注册
			// 创建UserProcess的实例
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)

		case 3:
			fmt.Println("=>>退出系统")
			os.Exit(0)
		default:
			fmt.Println("Sys>>您的输入有误，请重新输入")
		}
	}
}
