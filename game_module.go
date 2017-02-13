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
    //todo 全局对象 用singleton 利用golang once
    timerQueue = new(timer.TimerQueue)
    serviceMgr = new(service.ServiceMgr)

    timerQueue.OnInit()
    // test timer
    timerQueue.AddTimer(time.Second * 3, func() {
        fmt.Println("timer is out")
    })
    serviceMgr.OnInit()
    return true
}

func (gm *GameModule) Run(interval time.Duration, closeSig chan os.Signal) {
    ticker := time.NewTicker(interval)
    for {
        timerQueue.Select()
        serviceMgr.Select()
        //todo 网络部分netWork.Select()

        select {
        case <-ticker.C:
            //test ticker
            //todo 主逻辑中的tick驱动 npc ai...
            //example sceneMgr.OnTick() playerMgr.OnTick()
            fmt.Println("ticker")
        case <-closeSig:
            fmt.Println("closeSig")
            return
        default:
        }
    }
}

func (gm *GameModule) OnFinal() {
    timerQueue.OnFinal()
}
