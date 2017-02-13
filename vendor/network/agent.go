package network

import (
    "net"
    "log"
)

type Agent struct {
    uid     uint64
    conn    net.Conn //tcp or kcp
    ip      string
    svrId uint32 // gamesvr id
}

func newAgent(uid uint64, conn net.Conn) *Agent {
    return &Agent{uid : uid, conn, conn}
}

// 读取 解析 转发
func (a *Agent)Run() {
    for {
        data, err := a.conn.Read()
        if err != nil {
            log.Debug("read message: %v", err)
            break
        }

        if a.gate.Processor != nil {
            msg, err := a.gate.Processor.Unmarshal(data)
            if err != nil {
                log.Debug("unmarshal message error: %v", err)
                break
            }
            err = a.gate.Processor.Route(msg, a)
            if err != nil {
                log.Debug("route message error: %v", err)
                break
            }
        }
    }
}
