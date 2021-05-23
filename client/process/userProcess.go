package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_socket/client/message"
	"go_socket/client/model"
	"go_socket/client/utils"
	"net"
	"os"
)

type UserProcess struct {
	// 字段
}

// Login 完成登录
func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	// 1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial failed, err:", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType

	// 3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4.将loginMes序列化
	data, err := json.Marshal(loginMes) // data是[]byte类型
	if err != nil {
		fmt.Println("json.Marshal loginMes failed, err:", err)
		return
	}

	// 5.把data赋给mes.Data字段
	mes.Data = string(data)

	// 6.将mes进行序列化
	data, err = json.Marshal(mes) // 此时的data就是我们要发送的消息
	if err != nil {
		fmt.Println("json.Marshal mes failed, err:", err)
		return
	}

	// 7.先发送data的长度给服务器
	// 先获取到data的长度，然后转成[]byte
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(buf[0:4])
	if uint32(n) != 4 || err != nil {
		fmt.Println("conn.Write pkgLen slice failed, err:", err)
		return
	}
	//fmt.Println("客户端发送消息长度成功！，长度为", pkgLen)

	// 8.发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write data failed, err:", err)
		return
	}

	// 9.这里还需要处理服务器端返回的消息
	// 创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg failed, err:", err)
	}
	//fmt.Println("【Debug】独到的mes为:", mes)
	// 将 mes 的 Data 部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	//fmt.Println("【Debug】mes.Data:", mes.Data)
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	//fmt.Println("【Debug】解序列化后mes.Data:", mes.Data)
	if err != nil {
		fmt.Println("json.Unmarshal mes.Data failed, err:", err)
	}
	if loginResMes.Code == 200 {
		// 初始化CurUser的全局实例
		model.CurrentUser.Conn = conn
		model.CurrentUser.UserId = userId
		model.CurrentUser.UserStatus = message.UserOnline
		model.CurrentUser.UserName = loginResMes.UserName

		// 可以显示当前用户的在线列表遍历loginResMes.UsersId切片
		//fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UsersId {
			// 要求不显示自己id
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			// 完成客户端的onlineUsers的初始化
			user := &model.User{
				UserId:     v,
				UserName:   loginResMes.UserName,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println()

		// 这里我们还要在客户端启动一个协程，用来保持和服务器端的通讯
		// 如果服务器有数据推送给客户端，则接收并显示在客户端的终端
		go processServerMes(conn)

		// 2.循环显示我们登录成功的菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(">>", loginResMes.Error, "。")
	}
	return
}

// Register 完成注册
func (this *UserProcess) Register(userId int, userPwd, userName string) (err error) {
	// 1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial failed, err:", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType

	// 3.创建一个RegisterMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4.将RegisterMes序列化
	data, err := json.Marshal(registerMes) // data是[]byte类型
	if err != nil {
		fmt.Println("json.Marshal loginMes failed, err:", err)
		return
	}

	// 5.把data赋给mes.Data字段
	mes.Data = string(data)

	// 6.将mes进行序列化
	data, err = json.Marshal(mes) // 此时的data就是我们要发送的消息
	if err != nil {
		fmt.Println("json.Marshal mes failed, err:", err)
		return
	}

	// 创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	// 发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误, err:", err)
		return
	}

	// 从服务器端读取数据
	mes, err = tf.ReadPkg() // mes就是RegisterResMes
	if err != nil {
		fmt.Println("readPkg failed, err:", err)
	}
	// 将 mes 的 Data 部分反序列化成 RegisterResMes
	var registerResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Println("json.Unmarshal mes.Data failed, err:", err)
	}
	if registerResMes.Code == 200 {
		fmt.Println("=>>注册成功，请重新登录。")
	} else {
		fmt.Println(">>", registerResMes.Error, "。")
	}
	os.Exit(0)
	return
}
