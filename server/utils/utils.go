package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_socket/server/message"
	"net"
)

// 这里将这些方法关联到结构体中
type Transfer struct {
	// 分析需要的字段
	Conn net.Conn
	Buf  [8096]byte // 传输时使用的缓冲，数组
}

// readPkg
func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//buf := make([]byte, 8096)
	fmt.Println("等待读取客户端发送的数据...")
	n, err := this.Conn.Read(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Read header to buf failed, err:", err)
		return
	}

	// 根据buf[:4]转成一个uint32类型
	var pkgLen uint32 = binary.BigEndian.Uint32(this.Buf[0:4])
	// 根据pkgLen读取消息内容
	n, err = this.Conn.Read(this.Buf[:pkgLen]) //从套接字conn中读到buf中去
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.Read body to buffer failed, err:", err)
		return
	}
	// 把pkgLen反序列化成message.Message类型
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal failed, err", err)
		return
	}
	return
}

// writePkg
func (this *Transfer) WritePkg(data []byte) (err error) {

	// 1.先发送一个长度给对方
	// 先获取到data的长度，然后转成[]byte
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 发送长度
	n, err := this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write pkgLen slice failed, err:", err)
		return
	}
	fmt.Println("客户端发送消息长度成功！，长度为", pkgLen)

	// 2.发送data本身
	_, err = this.Conn.Write(data)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.Write data failed, err:", err)
		return
	}
	return
}
