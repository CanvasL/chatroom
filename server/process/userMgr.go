package process

import "fmt"

// 因为UserMgr实例在服务器端有且只有一个
//还因为在很多地方都会用到，所以我们将其定义为全局变量

var userMgr *UserMgr

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// InitUserMgr 完成对UserMgr的初始化工作
func InitUserMgr() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 完成对onlineUsers的添加/修改
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// DelOnlineUser 完成对onlineUsers的删除
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

// GetAllOnlineUser 完成对onlineUsers的查询，返回当前所有的在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// GetOnlineUserById 根据userId返回对应的在线用户
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	// 如何从map中取出一个值(待检测的方式)
	up, ok := this.onlineUsers[userId]
	if !ok {
		// 说明要查找的这个用户当前不在线
		err = fmt.Errorf("用户%d不存在", userId)
		return
	}
	return
}
