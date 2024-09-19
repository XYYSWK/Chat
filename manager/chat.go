package manager

import (
	socketio "github.com/googollee/go-socket.io"
	"sync"
	"time"
)

/*
每个账户各个客户端消息的发送
*/

const DefaultClientTimeout = time.Minute * 20

type ChatMap struct {
	// （GO 内置的 map 不是并发安全的，sync.Map 是并发安全的）
	m   sync.Map // k: accountID v: ConnMap（说明 accountID 可以不止有一个客户端设备）
	sID sync.Map // k: sID v: accountID，用于快速查找 某个连接 属于哪个 账户
}

type ConnMap struct {
	m sync.Map // k: sID v: ActiveConn（即每个 sID 对应一个活跃连接）
}

type ActiveConn struct {
	s          socketio.Conn // 连接对象，用于实际的连接操作（如发送/接收数据）
	activeTime time.Time     // 最后活动时间，用于判断连接是否超时
}

func NewChatMap() *ChatMap {
	return &ChatMap{m: sync.Map{}}
}

// Link 添加设备
func (c *ChatMap) Link(s socketio.Conn, accountID int64) {
	c.sID.Store(s.ID(), accountID) // 存入 SID 和 accountID 的对应关系
	cm, ok := c.m.Load(accountID)
	if !ok { // 没有找到对应的 ConnMap 对象，创建一个新的
		cm := &ConnMap{}
		activeConn := &ActiveConn{}
		activeConn.s = s
		activeConn.activeTime = time.Now()
		cm.m.Store(s.ID(), activeConn) // 将新连接存储在 ConnMap 中
		c.m.Store(accountID, cm)       // 将 ConnMap 存储在 c.m 中，以 accountID 为键
		return
	}
	activeConn := &ActiveConn{}
	activeConn.s = s
	activeConn.activeTime = time.Now()
	cm.(*ConnMap).m.Store(s.ID(), activeConn) // 将新的连接存储在 ConnMap 对象中
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
	// key: sID
	// value: ActiveConn
	cm.(*ConnMap).m.Range(func(key, value any) bool {
		activeConn := value.(*ActiveConn)
		activeConn.activeTime = time.Now() // 每次有消息发送，就重新计时
		activeConn.s.Emit(event, args...)  // 向指定客户端发送信息
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
			activeConn := value.(*ActiveConn)
			activeConn.activeTime = time.Now() // 每次有消息发送，就重新计时
			activeConn.s.Emit(event, args...)  // 向指定客户端发送信息
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
		f(value.(*ActiveConn).s)
		return true
	})
}

// HasSID 判断 SID 是否已经存在
func (c *ChatMap) HasSID(sID string) bool {
	_, ok := c.sID.Load(sID)
	return ok
}

// CheckForEachAllMap 遍历所有连接，检查是否超时，并关闭超时的连接
func (c *ChatMap) CheckForEachAllMap() {
	// c.m 是一个并发安全的映射，遍历每个 k-v 键值对
	c.m.Range(func(key, value any) bool {
		// key 是 account，value 是 ConnMap
		value.(*ConnMap).m.Range(func(key1, value1 any) bool {
			activeTime := value1.(*ActiveConn).activeTime
			if time.Now().Sub(activeTime) > DefaultClientTimeout { // 如果超时了
				err := value1.(*ActiveConn).s.Close()
				if err != nil {
					return false
				}
			}
			return true
		})
		return true
	})
}

func (c *ChatMap) CheckIsOnConnection(accountID int64) bool {
	_, ok := c.m.Load(accountID)
	return ok
}
