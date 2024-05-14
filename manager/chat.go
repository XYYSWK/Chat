package manager

import (
	socketio "github.com/googollee/go-socket.io"
	"sync"
)

/*
每个账户各个客户端消息的发送
*/

type ChatMap struct {
	// （GO 内置的 map 不是并发安全的，sync.Map 是并发安全的）
	m   sync.Map // k: accountID v: ConnMap（说明 accountID 可以不止有一个客户端设备）
	sID sync.Map // k: sID v: accountID（sID ——> socket.Conn.ID）
}

type ConnMap struct {
	m sync.Map // k: sID v: socketio.Conn
}

func NewChatMap() *ChatMap {
	return &ChatMap{m: sync.Map{}}
}

// Link 添加设备
func (c *ChatMap) Link(s socketio.Conn, accountID int64) {
	c.sID.Store(s.ID(), accountID) // 存入 SID 和 accountID 的对应关系
	cm, ok := c.m.Load(accountID)
	if !ok {
		cm := &ConnMap{}
		cm.m.Store(s.ID(), s)
		c.m.Store(accountID, cm)
		return
	}
	cm.(*ConnMap).m.Store(s.ID(), s)
}

// Leave 去除设备
func (c *ChatMap) Leave(s socketio.Conn) {
	accountID, ok := c.sID.LoadAndDelete(s.ID())
	if !ok {
		return
	}
	cm, ok := c.m.Load(accountID)
	if !ok {
		return
	}
	cm.(*ConnMap).m.Delete(s.ID())
	length := 0
	cm.(*ConnMap).m.Range(func(key, value any) bool {
		length++
		return true
	})
	if length == 0 {
		c.m.Delete(accountID)
	}
}

// Send 给指定账号的全部设备推送消息
func (c *ChatMap) Send(accountID int64, event string, args ...interface{}) {
	cm, ok := c.m.Load(accountID)
	if !ok { // 该账号不存在
		return
	}
	cm.(*ConnMap).m.Range(func(key, value any) bool {
		value.(socketio.Conn).Emit(event, args...) // 向指定客户端发送信息
		return true
	})
}

// SendMany 给指定多个账号的全部设备推送消息
// 参数：账号列表，事件名，要发送的数据
func (c *ChatMap) SendMany(accountIDs []int64, event string, args ...interface{}) {
	for _, accountID := range accountIDs {
		cm, ok := c.m.Load(accountID)
		if !ok { // 不存在该 accountID
			return
		}
		cm.(*ConnMap).m.Range(func(key, value interface{}) bool { // 遍历所有键值对
			value.(socketio.Conn).Emit(event, args...) // 向指定客户端发送信息
			return true
		})
	}
}

// SendAll 给全部设备推送消息
func (c *ChatMap) SendAll(event string, args ...interface{}) {
	c.m.Range(func(key, value any) bool {
		value.(*ConnMap).m.Range(func(key, value any) bool {
			value.(socketio.Conn).Emit(event, args...)
			return true
		})
		return true
	})
}

type EachFunc socketio.EachFunc // 定义每个客户端连接的处理函数

// ForEach 遍历指定账号的全部设备
func (c *ChatMap) ForEach(accountID int64, f EachFunc) {
	cm, ok := c.m.Load(accountID)
	if !ok {
		return
	}
	cm.(*ConnMap).m.Range(func(key, value any) bool {
		f(value.(socketio.Conn))
		return true
	})
}

// HasSID 判断 SID 是否已经存在
func (c *ChatMap) HasSID(sID string) bool {
	_, ok := c.sID.Load(sID)
	return ok
}
