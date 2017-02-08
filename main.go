package main

import (
    "timer"
    "fmt"
    "time"
)

func main()  {
    timerQueue := new(timer.TimerQueue)
    fmt.Println("servicemgr start")
    timerQueue.AddTimer(time.Duration(time.Second + 3), func() {
        fmt.Println("hello world")
    })
    for {
        timerQueue.OnTick()
    }
}