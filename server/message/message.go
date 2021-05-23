package message

import "go_socket/server/model"

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

// 定义几个用户在线的常量
const (
	UserOnline  = "在线"
	UserOffline = "离线"
	UserBusy    = "忙碌"
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息内容
}

// 定义两个消息，后续可增加

type LoginMes struct {
	UserId   int    `json:"userId"`   //用户id
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code     int    `json:"code"` //返回状态码 500表示用户未注册；200表示登录成功
	UsersId  []int  // 增加字段，保存用户id的切片
	UserName string `json:"userName"`
	Error    string `json:"error"` //返回错误
}

type RegisterMes struct {
	User model.User `json:"user"`
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码 400表示用户已存在；200表示注册成功
	Error string `json:"error"` //返回错误
}

// 为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int    `json:"userId"`
	Status string `json:"status"`
}

// 增加一个SmsMes
type SmsMes struct {
	Content string `json:"content"` // 内容
	model.User
}
