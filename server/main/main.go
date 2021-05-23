package main

import (
	"fmt"
	"go_socket/server/model"
	"go_socket/server/process"
	"net"
	"time"
)

// initUserDao 完成对UserDao的初始化任务
//初始化顺序：要先初始化redis，再初始化UserDao
func initUserDao() {
	// 这里的pool本身就是一个全局的变量，在main包的redis.go中
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	// 服务器启动时，我们就去初始化redis的连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	// 初始化UserDao
	initUserDao()
	// 初始化Mar
	process.InitUserMgr()
	// 提示信息
	fmt.Println("服务器在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen failed, err:", err)
		return
	}
	defer listen.Close()
	// 一旦监听成功，就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端来连接服务器...")
		conn, err := listen.Accept() //得到套接字conn
		if err != nil {
			fmt.Println("listen.Accept failed, err:", err)
		}
		// 一旦连接成功，则启动协程和客户端保持通讯
		go commWithClient(conn)
	}
}

// commWithClient 处理和客户端的通讯
func commWithClient(conn net.Conn) {
	// 这里需要延时关闭conn
	defer conn.Close()

	// 调用总控
	// 创建总控实例
	processor := &Processor{
		Conn: conn,
	}
	err := processor.processMes()
	if err != nil {
		fmt.Println("客户端和服务器端通信的协程错误，err:", err)
		return
	}
}
