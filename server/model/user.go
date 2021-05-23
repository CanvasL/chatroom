package model

// 定义一个用户的结构体
type User struct {
	// 确定字段信息
	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus string `json:"userStatus"` //用户的状态..
}
