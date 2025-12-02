package kcp

import (
	"github.com/goodluck0107/gcore/gerrors"
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/gnetwork"
	"github.com/goodluck0107/gcore/gpacket"
	"github.com/goodluck0107/gcore/gutils/gcall"
	"github.com/goodluck0107/gcore/gutils/gnet"
	"github.com/goodluck0107/gcore/gutils/gtime"
	"github.com/xtaci/kcp-go/v5"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type serverConn struct {
	rw                sync.RWMutex    // 锁
	id                int64           // 连接ID
	uid               int64           // 用户ID
	state             int32           // 连接状态
	conn              *kcp.UDPSession // UDP源连接
	connMgr           *serverConnMgr  // 连接管理
	chWrite           chan chWrite    // 写入队列
	lastHeartbeatTime int64           // 上次心跳时间
	done              chan struct{}   // 写入完成信号
	close             chan struct{}   // 关闭信号
}

var _ gnetwork.Conn = &serverConn{}

// ID 获取连接ID
func (c *serverConn) ID() int64 {
	return c.id
}

// UID 获取用户ID
func (c *serverConn) UID() int64 {
	return atomic.LoadInt64(&c.uid)
}

// Bind 绑定用户ID
func (c *serverConn) Bind(uid int64) {
	atomic.StoreInt64(&c.uid, uid)
}

// Unbind 解绑用户ID
func (c *serverConn) Unbind() {
	atomic.StoreInt64(&c.uid, 0)
}

// Send 发送消息（同步）
func (c *serverConn) Send(msg []byte) (err error) {
	if err = c.checkState(); err != nil {
		return
	}

	c.rw.RLock()
	conn := c.conn
	c.rw.RUnlock()

	if conn == nil {
		return gerrors.ErrConnectionClosed
	}

	_, err = conn.Write(msg)
	return
}

// Push 发送消息（异步）
func (c *serverConn) Push(msg []byte) (err error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err = c.checkState(); err != nil {
		return
	}

	c.chWrite <- chWrite{typ: dataPacket, msg: msg}

	return
}

// State 获取连接状态
func (c *serverConn) State() gnetwork.ConnState {
	return gnetwork.ConnState(atomic.LoadInt32(&c.state))
}

// Close 关闭连接
func (c *serverConn) Close(force ...bool) error {
	if len(force) > 0 && force[0] {
		return c.forceClose(true)
	} else {
		return c.graceClose(true)
	}
}

// LocalIP 获取本地IP
func (c *serverConn) LocalIP() (string, error) {
	addr, err := c.LocalAddr()
	if err != nil {
		return "", err
	}

	return gnet.ExtractIP(addr)
}

// LocalAddr 获取本地地址
func (c *serverConn) LocalAddr() (net.Addr, error) {
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
func (c *serverConn) RemoteIP() (string, error) {
	addr, err := c.RemoteAddr()
	if err != nil {
		return "", err
	}

	return gnet.ExtractIP(addr)
}

// RemoteAddr 获取远端地址
func (c *serverConn) RemoteAddr() (net.Addr, error) {
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

// 初始化连接
func (c *serverConn) init(cm *serverConnMgr, id int64, conn *kcp.UDPSession) {
	c.id = id
	c.conn = conn
	c.connMgr = cm
	c.chWrite = make(chan chWrite, 4096)
	c.done = make(chan struct{})
	c.close = make(chan struct{})
	c.lastHeartbeatTime = gtime.Now().UnixNano()
	atomic.StoreInt64(&c.uid, 0)
	atomic.StoreInt32(&c.state, int32(gnetwork.ConnOpened))

	conn.SetStreamMode(true)
	conn.SetWriteDelay(false)
	conn.SetNoDelay(1, 10, 2, 1)
	conn.SetMtu(1500)
	//conn.SetWindowSize(config.SndWnd, config.RcvWnd)
	conn.SetACKNoDelay(true)

	gcall.Go(c.read)

	gcall.Go(c.write)

	if c.connMgr.server.connectHandler != nil {
		c.connMgr.server.connectHandler(c)
	}
}

// 检测连接状态
func (c *serverConn) checkState() error {
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
func (c *serverConn) graceClose(isNeedRecycle bool) error {
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

	if isNeedRecycle {
		c.connMgr.recycle(conn)
	}

	if c.connMgr.server.disconnectHandler != nil {
		c.connMgr.server.disconnectHandler(c)
	}

	return err
}

// 强制关闭
func (c *serverConn) forceClose(isNeedRecycle bool) error {
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

	if isNeedRecycle {
		c.connMgr.recycle(conn)
	}

	if c.connMgr.server.disconnectHandler != nil {
		c.connMgr.server.disconnectHandler(c)
	}

	return err
}

// 读取消息
func (c *serverConn) read() {
	conn := c.conn

	for {
		select {
		case <-c.close:
			return
		default:
			msg, err := gpacket.ReadMessage(conn)
			if err != nil {
				_ = c.forceClose(true)
				return
			}

			if c.connMgr.server.opts.heartbeatInterval > 0 {
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
				// responsive heartbeat
				if c.connMgr.server.opts.heartbeatMechanism == RespHeartbeat {
					if heartbeat, err := gpacket.PackHeartbeat(); err != nil {
						glog.Errorf("pack heartbeat message error: %v", err)
					} else {
						if _, err = conn.Write(heartbeat); err != nil {
							glog.Errorf("write heartbeat message error: %v", err)
						}
					}
				}
				continue
			}

			// ignore empty packet
			if len(msg) == 0 {
				continue
			}

			if c.connMgr.server.receiveHandler != nil {
				c.connMgr.server.receiveHandler(c, msg)
			}
		}
	}
}

// 写入消息
func (c *serverConn) write() {
	var (
		conn   = c.conn
		ticker *time.Ticker
	)

	if c.connMgr.server.opts.heartbeatInterval > 0 {
		ticker = time.NewTicker(c.connMgr.server.opts.heartbeatInterval)
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
			deadline := gtime.Now().Add(-2 * c.connMgr.server.opts.heartbeatInterval).UnixNano()
			if atomic.LoadInt64(&c.lastHeartbeatTime) < deadline {
				glog.Debugf("connection heartbeat timeout, cid: %d", c.id)
				_ = c.forceClose(true)
				return
			} else {
				if c.connMgr.server.opts.heartbeatMechanism == TickHeartbeat {
					if c.isClosed() {
						return
					}

					if heartbeat, err := gpacket.PackHeartbeat(); err != nil {
						glog.Errorf("pack heartbeat message error: %v", err)
					} else {
						// send heartbeat packet
						if _, err = conn.Write(heartbeat); err != nil {
							glog.Errorf("write heartbeat message error: %v", err)
						}
					}
				}
			}
		}
	}
}

// 是否已关闭
func (c *serverConn) isClosed() bool {
	return gnetwork.ConnState(atomic.LoadInt32(&c.state)) == gnetwork.ConnClosed
}
