package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// 服务器启动时就要创建UserDao的实例，把它声明成全局变量
//在需要和redis操作时，就直接使用即可
var (
	MyUserDao *UserDao
)

// 定义一个UserDao结构体，完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// NewUserDao 使用工厂模式创建一个UserDao的实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// getUserById 根据用户id返回一个User示例
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	// 通过给定的id去redis中查询用户
	res, err := redis.String(conn.Do("HGET", "users", id))
	if err != nil {
		// 错误
		if err == redis.ErrNil {
			// 表示在users这个hash中没有找到对应的id
			err = ERROR_USER_NOTEXSITS
		}
		return
	}

	// 执行到这里，说明users中存在该用户
	// 这里我们需要把res反序列化成一个User对象
	user = &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal res failed, err:", err)
		return
	}
	return
}

// Login 根据提供的用户id和密码完成登录的校验
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	// 先从UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		// 错误一：用户不存在
		return
	}

	// 此时用户user已获取到,此user已被反序列化过了
	if user.UserPwd != userPwd {
		// 错误二：密码不正确
		err = ERROR_USER_PWD
		return
	}

	return
}

// Register 根据提供的用户id、密码和昵称完成注册
func (this *UserDao) Register(user *User) (err error) {
	// 先从UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		// 错误一：用户已存在
		err = ERROR_USER_EXISTS
		return
	}
	// 此时用户id在redis中还没有注册过
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	// 入库
	_, err = conn.Do("HSET", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("注册用户信息入库错误, err:", err)
		return
	}
	return
}
