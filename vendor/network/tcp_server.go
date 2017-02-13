package network

import (
    "net"
    "log"
    "os"
    "fmt"
)

type TcpServer struct {
    Addr    string
    MaxConnNum int
    ln net.TCPListener
    PendingWriteNum int
}

func (t *TcpServer) Start() {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", 8001)
    if err != nil {
        fmt.Println(err)
    }

    ln, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        fmt.Println(err)
    }
    t.ln = ln

    for {
        conn, err := t.ln.AcceptTCP()
        if err != nil {
            fmt.Println(err)
        }
        tcpConn := newTcpConn(conn, t.PendingWriteNum)
        agent := newAgent(1, tcpConn)
        go func() {
            agent.Run()
        }()
    }

}
