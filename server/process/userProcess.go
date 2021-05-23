package process

import (
	"encoding/json"
	"fmt"
	"go_socket/server/message"
	"go_socket/server/model"
	"go_socket/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	// 增加字段，表示该Conn是哪个用户的
	UserId int
}

// NotifyOthersOnlineUser 编写通知所有在线的用户的方法
//userId要通知其他的在线用户，我上线
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	// 遍历 onlineUsers，然后一个一个地发送
	for id, up := range userMgr.onlineUsers {
		// 过滤掉自己
		if id == userId {
			continue
		}
		// 开始通知 单独写一个方法
		up.NotifyOnlineUsers(userId)
	}
}

func (this *UserProcess) NotifyOnlineUsers(userId int) {
	// 组装我们的消息
	mes := message.Message{
		Type: message.NotifyUserStatusMesType,
	}
	notifyUserStatusMes := message.NotifyUserStatusMes{
		UserId: userId,
		Status: message.UserOnline,
	}

	// 将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal notifyUserStatusMes failed, err:", err)
		return
	}
	// 将序列化后的notifyUserStatusMes赋给Data
	mes.Data = string(data)
	// 对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal mes failed, err:", err)
		return
	}

	// 发送，创建一个Transfer的实例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyOnlineUsers failed, err:", err)
	}
}

// ServerProcessLogin 专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1.先从 message 中取出 mes.Data , 并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal mes.Data failed, err:", err)
		return
	}

	// 2.先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	// 3.再声明一个 LoginResMes，并完成赋值
	var loginResMes message.LoginResMes

	// 使用UserDao实例去redis数据库完成登录验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		// 不合法
		if err == model.ERROR_USER_NOTEXSITS {
			// 错误一：用户不存在
			loginResMes.Code = 500 // 500这个状态码表示该用户不存在
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			// 错误二：密码有问题
			loginResMes.Code = 403 // 403这个状态码表示密码错误
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505 // 未知错误
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		// 合法
		//fmt.Println("【Debug】可以登录code=200,loginResMes=", loginResMes)
		loginResMes.Code = 200
		loginResMes.UserName = user.UserName
		// 因为用户登录成功，我们就把该登录成功的用户放入到userMgr中
		// 将登录成功的用户id赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		// 通知其他在线用户此用户的登录信息
		this.NotifyOthersOnlineUser(this.UserId)
		// 将当前在线用户的id放入到loginResMes.User中去
		// 遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println("=>>", user.UserName, "登录成功。")
	}
	//fmt.Println("【Debug】loginResMes=", loginResMes)

	// 4.将 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal loginResMes failed, err:", err)
		return
	}

	// 5.将 data 赋值给 resMes 结构体
	resMes.Data = string(data)
	//fmt.Println("【Debug】resMes.Data:", resMes.Data)

	// 6.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal resMes failed, err:", err)
		return
	}
	fmt.Println("【Debug】发送的data为:", resMes)

	// 7.发送 data ，我们将其封装到 writePkg 函数中
	// 因为使用了分层模式，我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// ServerProcessRegister 专门处理注册请求
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 1.先从 message 中取出 mes.Data , 并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal mes.Data failed, err:", err)
		return
	}

	// 2.先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	// 3.再声明一个 RegisterResMes，并完成赋值
	var registerResMes message.RegisterResMes

	// 使用UserDao实例去redis数据库完成注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		// 错误
		if err == model.ERROR_USER_EXISTS {
			// 错误一：用户已存在
			registerResMes.Code = 505 // 500这个状态码表示该用户不存在
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 506 // 未知错误
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		// 合法
		registerResMes.Code = 200
		fmt.Println("=>>id为", registerMes.User.UserId, "、昵称为", registerMes.User.UserName, "的用户注册成功。")
	}

	// 4.将 loginResMes 序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal loginResMes failed, err:", err)
		return
	}

	// 5.将 data 赋值给 resMes 结构体
	resMes.Data = string(data)

	// 6.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal resMes failed, err:", err)
		return
	}

	// 7.发送 data ，我们将其封装到 writePkg 函数中
	// 因为使用了分层模式，我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
