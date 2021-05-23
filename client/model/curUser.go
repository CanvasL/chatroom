package model

import "net"

// 因为在客户端我们很多地方会使用到curUser，所以做成全局变量
// 在用户登录成功后就要完成对curUser的初始化
var CurrentUser CurUser

type CurUser struct {
	Conn net.Conn
	User
}
