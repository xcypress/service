##服务框架 vendor/service
* 包管理 govender
* 基于docker linux 3.10
* etcd
* registrator
* etcd-browser

##功能
* 支持自动服务的注册和发现
* 支持服务负载均衡
* CP 一致性和分区容忍性
* 支持http 和 tcp
* 安全：TLS
* 可靠: etcd使用Raft适当分布

###Usage
```
func main()  {
    serviceMgr := new(service.ServiceMgr)
    serviceMgr.OnInit()
    fmt.Println("servicemgr start")
    for {
        serviceMgr.OnTick()
    }
}
```

###Appliction 模块
* 基于Ioc 控制反转
* 游戏逻辑单goroutine 开发简单高效
* 类似Half Sync/Half Async net timer ticker单独跑goroutine --通过chan mq--> main goroutine
* 游戏逻辑靠客户端Req和服务器Tick驱动
* TimerQueue timer ticker trans to main loop
###Usage
```
type Module interface {
    OnInit() bool
    OnFinal()
    Run(interval time.Duration, closeSig chan os.Signal)
}

examle：
GameModule 实现Module接口
将GameModule 注册到app
调用app.Run() 程序即进入MainLoop

func main() {
    app := &module.Application{}
    GameModule := &GameModule{}

    app.Register(GameModule)
    app.SetInterval(time.Second * DefaultInterval)
    app.Run()
}
```