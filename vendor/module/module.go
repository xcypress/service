package module

import (
    "os"
    "os/signal"
)

//Ioc skeleton
type Module interface {
    OnInit() bool
    OnFinal()
    Run(closeSig chan os.Signal)
}

type Application struct {
    module Module
}

func (a *Application) Register(app Module) {
    a.module = app
}


func (a *Application) Run() {
    if a.module.OnInit() == false {
        a.module.OnFinal()
        return
    }

    closeSig := make(chan os.Signal, 1)
    signal.Notify(closeSig, os.Interrupt, os.Kill)

    a.module.Run(closeSig)

    a.module.OnFinal()
}
