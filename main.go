package main

import (
    "module"
    "time"
)

const DefaultInterval = 5

func main() {
    app := &module.Application{}
    GameModule := &GameModule{}

    app.Register(GameModule)
    app.SetInterval(time.Second * DefaultInterval)
    app.Run()
}