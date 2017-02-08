package module

import (
    "time"
    "os"
    "os/signal"
)

//Ioc skeleton
type Module interface {
    OnInit() bool
    OnFinal()
    Run(interval time.Duration, closeSig chan os.Signal)
}

type Application struct {
    module Module
    interval time.Duration
}

func (a *Application) Register(app Module) {
    a.module = app
}

func (a *Application) SetInterval(interval time.Duration) {
    a.interval = interval
}

func (a *Application) Run() {
    if a.module.OnInit() == false {
        a.module.OnFinal()
        return
    }

    closeSig := make(chan os.Signal, 1)
    signal.Notify(closeSig, os.Interrupt, os.Kill)

    a.module.Run(a.interval, closeSig)

    a.module.OnFinal()
}
