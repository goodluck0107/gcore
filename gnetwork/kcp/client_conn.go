package kcp

import (
	"gitee.com/monobytes/gcore/gerrors"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gnetwork"
	"gitee.com/monobytes/gcore/gpacket"
	"gitee.com/monobytes/gcore/gutils/gcall"
	"gitee.com/monobytes/gcore/gutils/gnet"
	"gitee.com/monobytes/gcore/gutils/gtime"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type clientConn struct {
	rw                sync.RWMutex
	id                int64         // 连接ID
	uid               int64         // 用户ID
	conn              net.Conn      // TCP源连接
	state             int32         // 连接状态
	client            *client       // 客户端
	chWrite           chan chWrite  // 写入队列
	done              chan struct{} // 写入完成信号
	close             chan struct{} // 关闭信号
	lastHeartbeatTime int64         // 上次心跳时间
}

var _ gnetwork.Conn = &clientConn{}

func newClientConn(client *client, id int64, conn net.Conn) gnetwork.Conn {
	c := &clientConn{
		id:                id,
		conn:              conn,
		state:             int32(gnetwork.ConnOpened),
		client:            client,
		chWrite:           make(chan chWrite, 4096),
		done:              make(chan struct{}),
		close:             make(chan struct{}),
		lastHeartbeatTime: gtime.Now().UnixNano(),
	}

	gcall.Go(c.read)

	gcall.Go(c.write)

	if c.client.connectHandler != nil {
		c.client.connectHandler(c)
	}

	return c
}

// ID 获取连接ID
func (c *clientConn) ID() int64 {
	return c.id
}

// UID 获取用户ID
func (c *clientConn) UID() int64 {
	return atomic.LoadInt64(&c.uid)
}

// Bind 绑定用户ID
func (c *clientConn) Bind(uid int64) {
	atomic.StoreInt64(&c.uid, uid)
}

// Unbind 解绑用户ID
func (c *clientConn) Unbind() {
	atomic.StoreInt64(&c.uid, 0)
}

// Send 发送消息（同步）
func (c *clientConn) Send(msg []byte) error {
	if err := c.checkState(); err != nil {
		return err
	}

	c.rw.RLock()
	conn := c.conn
	c.rw.RUnlock()

	if conn == nil {
		return gerrors.ErrConnectionClosed
	}

	_, err := conn.Write(msg)
	return err
}

// Push 发送消息（异步）
func (c *clientConn) Push(msg []byte) (err error) {
	if err = c.checkState(); err != nil {
		return
	}

	c.rw.RLock()
	c.chWrite <- chWrite{typ: dataPacket, msg: msg}
	c.rw.RUnlock()

	return
}

// State 获取连接状态
func (c *clientConn) State() gnetwork.ConnState {
	return gnetwork.ConnState(atomic.LoadInt32(&c.state))
}

// Close 关闭连接
func (c *clientConn) Close(force ...bool) error {
	if len(force) > 0 && force[0] {
		return c.forceClose()
	} else {
		return c.graceClose()
	}
}

// LocalIP 获取本地IP
func (c *clientConn) LocalIP() (string, error) {
	addr, err := c.LocalAddr()
	if err != nil {
		return "", err
	}

	return gnet.ExtractIP(addr)
}

// LocalAddr 获取本地地址
func (c *clientConn) LocalAddr() (net.Addr, error) {
	if err := c.checkState(); err != nil {
		return nil, err
	}

	c.rw.RLock()
	conn := c.conn
	c.rw.RUnlock()

	if conn == nil {
		return nil, gerrors.ErrConnectionClosed
	}

	return conn.LocalAddr(), nil
}

// RemoteIP 获取远端IP
func (c *clientConn) RemoteIP() (string, error) {
	addr, err := c.RemoteAddr()
	if err != nil {
		return "", err
	}

	return gnet.ExtractIP(addr)
}

// RemoteAddr 获取远端地址
func (c *clientConn) RemoteAddr() (net.Addr, error) {
	if err := c.checkState(); err != nil {
		return nil, err
	}

	c.rw.RLock()
	conn := c.conn
	c.rw.RUnlock()

	if conn == nil {
		return nil, gerrors.ErrConnectionClosed
	}

	return conn.RemoteAddr(), nil
}

// 检测连接状态
func (c *clientConn) checkState() error {
	switch gnetwork.ConnState(atomic.LoadInt32(&c.state)) {
	case gnetwork.ConnHanged:
		return gerrors.ErrConnectionHanged
	case gnetwork.ConnClosed:
		return gerrors.ErrConnectionClosed
	default:
		return nil
	}
}

// 优雅关闭
func (c *clientConn) graceClose() error {
	if !atomic.CompareAndSwapInt32(&c.state, int32(gnetwork.ConnOpened), int32(gnetwork.ConnHanged)) {
		return gerrors.ErrConnectionNotOpened
	}

	c.rw.RLock()
	c.chWrite <- chWrite{typ: closeSig}
	c.rw.RUnlock()

	<-c.done

	if !atomic.CompareAndSwapInt32(&c.state, int32(gnetwork.ConnHanged), int32(gnetwork.ConnClosed)) {
		return gerrors.ErrConnectionNotHanged
	}

	c.rw.Lock()
	close(c.chWrite)
	close(c.close)
	close(c.done)
	conn := c.conn
	c.conn = nil
	c.rw.Unlock()

	err := conn.Close()

	if c.client.disconnectHandler != nil {
		c.client.disconnectHandler(c)
	}

	return err
}

// 强制关闭
func (c *clientConn) forceClose() error {
	if !atomic.CompareAndSwapInt32(&c.state, int32(gnetwork.ConnOpened), int32(gnetwork.ConnClosed)) {
		if !atomic.CompareAndSwapInt32(&c.state, int32(gnetwork.ConnHanged), int32(gnetwork.ConnClosed)) {
			return gerrors.ErrConnectionClosed
		}
	}

	c.rw.Lock()
	close(c.chWrite)
	close(c.close)
	close(c.done)
	conn := c.conn
	c.conn = nil
	c.rw.Unlock()

	err := conn.Close()

	if c.client.disconnectHandler != nil {
		c.client.disconnectHandler(c)
	}

	return err
}

// 读取消息
func (c *clientConn) read() {
	conn := c.conn

	for {
		select {
		case <-c.close:
			return
		default:
			msg, err := gpacket.ReadMessage(conn)
			if err != nil {
				_ = c.forceClose()
				return
			}

			if c.client.opts.heartbeatInterval > 0 {
				atomic.StoreInt64(&c.lastHeartbeatTime, gtime.Now().UnixNano())
			}

			switch c.State() {
			case gnetwork.ConnHanged:
				continue
			case gnetwork.ConnClosed:
				return
			default:
				// ignore
			}

			isHeartbeat, err := gpacket.CheckHeartbeat(msg)
			if err != nil {
				glog.Errorf("check heartbeat message error: %v", err)
				continue
			}

			// ignore heartbeat packet
			if isHeartbeat {
				continue
			}

			// ignore empty packet
			if len(msg) == 0 {
				continue
			}

			if c.client.receiveHandler != nil {
				c.client.receiveHandler(c, msg)
			}
		}
	}
}

// 写入消息
func (c *clientConn) write() {
	var (
		conn   = c.conn
		ticker *time.Ticker
	)

	if c.client.opts.heartbeatInterval > 0 {
		ticker = time.NewTicker(c.client.opts.heartbeatInterval)
		defer ticker.Stop()
	} else {
		ticker = &time.Ticker{C: make(chan time.Time, 1)}
	}

	for {
		select {
		case r, ok := <-c.chWrite:
			if !ok {
				return
			}

			if r.typ == closeSig {
				c.rw.RLock()
				c.done <- struct{}{}
				c.rw.RUnlock()
				return
			}

			if c.isClosed() {
				return
			}

			if _, err := conn.Write(r.msg); err != nil {
				glog.Errorf("write data message error: %v", err)
			}
		case <-ticker.C:
			deadline := gtime.Now().Add(-2 * c.client.opts.heartbeatInterval).UnixNano()
			if atomic.LoadInt64(&c.lastHeartbeatTime) < deadline {
				glog.Debugf("connection heartbeat timeout")
				_ = c.forceClose()
				return
			} else {
				if c.isClosed() {
					return
				}

				if heartbeat, err := gpacket.PackHeartbeat(); err != nil {
					glog.Errorf("pack heartbeat message error: %v", err)
				} else {
					// send heartbeat packet
					if _, err := conn.Write(heartbeat); err != nil {
						glog.Errorf("write heartbeat message error: %v", err)
					}
				}
			}
		}
	}
}

// 是否已关闭
func (c *clientConn) isClosed() bool {
	return gnetwork.ConnState(atomic.LoadInt32(&c.state)) == gnetwork.ConnClosed
}
