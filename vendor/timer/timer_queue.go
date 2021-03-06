package timer

import (
    "time"
    "fmt"
)

type  SimpleTimer struct {
    t *time.Timer
    cb func()
}

type SimpleTicker struct {
    t *time.Ticker
    cb func()
}

func (s *SimpleTicker) Stop() {
    s.t.Stop()
    s.cb = nil
}

func (s *SimpleTimer) Stop() {
    s.t.Stop()
    s.cb = nil
}

type TimerQueue struct {
    TimerMQ chan *SimpleTimer
    TickerMQ chan *SimpleTicker
}

func (tq *TimerQueue) OnInit() {
    tq.TimerMQ = make(chan *SimpleTimer, 100)
    tq.TickerMQ = make(chan *SimpleTicker, 100)
}

func (tq *TimerQueue) Select() {
    select {
    case timer := <- tq.TimerMQ:
        timer.cb()
    case ticker := <- tq.TickerMQ:
        ticker.cb()
    default:
    }
}

func (tq *TimerQueue) OnFinal() {
    for timer := range tq.TimerMQ {
        timer.Stop()
    }
    for ticker := range tq.TimerMQ {
        ticker.Stop()
    }
    fmt.Println("timerQueue onfinal")
}

func (tq *TimerQueue) AddTimer(d time.Duration, f func()) *SimpleTimer {
    timer := new(SimpleTimer)
    timer.cb = f
    timer.t = time.AfterFunc(d, func() {
        tq.TimerMQ <- timer
    })
    return timer
}

func (tq *TimerQueue) AddTicker(d time.Duration, f func()) *SimpleTicker{
    ticker := new(SimpleTicker)
    ticker.cb = f
    ticker.t = time.NewTicker(d)

    for {
        select {
        case <-ticker.t.C:
            tq.TickerMQ <- ticker
        }
    }
    return ticker
}



