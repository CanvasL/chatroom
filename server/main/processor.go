package main

import (
	"fmt"
	"go_socket/server/message"
	"go_socket/server/process"
	"go_socket/server/utils"
	"io"
	"net"
)

// 先创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

// serverProcessMes 根据客户端发送的消息种类不同，决定调用哪个函数来处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		// 处理登录
		// 创建一个UserProcess实例
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)

	case message.RegisterMesType:
		// 处理注册
		// 创建一个UserProcess实例
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)

	case message.SmsMesType:
		// 处理群发
		// 创建一个SmsProcess的实例，完成转发群聊消息的任务
		smsProcess := &process.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

func (this *Processor) processMes() (err error) {
	// 循环读客户端发送的信息
	for {
		// 这里我们将读取数据包，直接封装成一个函数readPkg(), 返回Message， Err
		// 创建一个Transfer的实例，完成读包的若任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端已退出。服务器端也退出...")
				return err
			}
			fmt.Println("readPkg conn failed, err:", err)
			return err
		}
		fmt.Println("mes=", mes)

		err = this.serverProcessMes(&mes)
		if err != nil {
			fmt.Println("serverProcessMes failed, err:", err)
			return err
		}
	}
}
