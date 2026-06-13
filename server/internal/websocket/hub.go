package websocket

import (
	"errors"
	"sync"
)

// 定义一个错误类型
var ErrTooManyConnections = errors.New("too many websocket connections")

// Hub 连接中心
type Hub struct {
	mu             sync.RWMutex
	connections    map[string]*Connection
	userConns      map[string]map[string]*Connection
	maxConnections int
}

// 初始化连接中心
func NewHub(maxConnections int) *Hub {
	if maxConnections <= 0 {
		maxConnections = 10000
	}
	// 初始化连接中心
	return &Hub{
		// 这是所有的连接，靠ID来查找
		connections: make(map[string]*Connection),
		// 这是所有用户的所有连接，靠用户ID来查找第一次，再靠连接ID来查找第二次
		userConns:      make(map[string]map[string]*Connection),
		maxConnections: maxConnections,
	}
}

// Register 注册连接
func (h *Hub) Register(conn *Connection) error {
	// 加锁，确保线程安全
	h.mu.Lock()
	defer h.mu.Unlock()

	// 检查是否超过最大连接数
	if len(h.connections) >= h.maxConnections {
		return ErrTooManyConnections
	}

	// 注册连接
	h.connections[conn.ID()] = conn
	// 判断该用户在userconns中是否存在
	if _, ok := h.userConns[conn.UserID()]; !ok {
		// 如果不存在，在userconns中创建该用户的连接map
		h.userConns[conn.UserID()] = make(map[string]*Connection)
	}
	// 注册连接到userconns中该用户的连接map
	h.userConns[conn.UserID()][conn.ID()] = conn
	return nil
}

// Unregister 注销连接
func (h *Hub) Unregister(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// map的专属操作，删除连接
	delete(h.connections, conn.ID())
	userID := conn.UserID()
	// 判断该用户在userconns中是否存在
	if conns, ok := h.userConns[userID]; ok {
		// 如果存在，在userconns中删除该用户的连接
		delete(conns, conn.ID())
		// 假如该用户的连接数为0，删除该用户的连接map
		if len(conns) == 0 {
			delete(h.userConns, userID)
		}
	}
}

// ActiveConnections 获取当前连接数
func (h *Hub) ActiveConnections() int {
	// 加一个读锁，确保线程安全
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.connections)
}

// RWMutex 读写锁的lock方法 和 Mutex 的lock方法效果一致但是RWMuutex的性能差一些
