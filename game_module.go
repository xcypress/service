package main

import (
    "timer"
    "service"
    "time"
    "os"
    "fmt"
)

var (
    timerQueue  *timer.TimerQueue
    serviceMgr  *service.ServiceMgr
)

type GameModule struct {

}

func (gm *GameModule) OnInit() bool {
    timerQueue.OnInit()
    timerQueue.AddTimer(time.Second*5, func() {
        fmt.Println("timer is out")
    })
    serviceMgr.OnInit()
    return true
}

func (gm *GameModule) Run(interval time.Duration, closeSig os.Signal) {
    ticker := time.NewTicker(interval)
    for {
        timerQueue.OnTick()
        serviceMgr.OnTick()

        select {
        case ticker.C:
            fmt.Println("ticker")
        case <-closeSig:
        //todo 结束循环
            return
        }
    }
}

func (gm *GameModule) OnFinal() {
    timerQueue.OnFinal()
}
