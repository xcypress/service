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
    timerQueue = new(timer.TimerQueue)
    serviceMgr = new(service.ServiceMgr)

    timerQueue.OnInit()
    timerQueue.AddTimer(time.Second*5, func() {
        fmt.Println("timer is out")
    })
    serviceMgr.OnInit()
    return true
}

func (gm *GameModule) Run(interval time.Duration, closeSig chan os.Signal) {
    fmt.Println(interval)
    ticker := time.NewTicker(time.Second*5)
    for {
        timerQueue.OnTick()
        serviceMgr.OnTick()

        select {
        case <-ticker.C:
            fmt.Println("ticker")
        case <-closeSig:
            return
        }
    }
}

func (gm *GameModule) OnFinal() {
    timerQueue.OnFinal()
}
