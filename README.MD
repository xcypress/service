## 服务框架 vendor/service
* govender
* docker linux 3.10
* etcd
* registrator
* etcd-browser

## 功能
* 支持自动服务的注册和发现
* 支持服务负载均衡 无状态服务器 比如gateway
* 满足cap中cp 一致性和分区容忍性
* 支持http 和 tcp
* 安全：TLS
* 可靠: etcd使用Raft适当分布

### 调研
* grpc 适合web和简单的service 支持req/rsp stream http2 复用 流控 但使用会破坏主逻辑单goroutine，鉴于mmo游戏逻辑复杂，尽量保证主逻辑单goroutine
* zookeeper 复杂难用
* consul 各种功能齐全 web ui漂亮 但是watch功能有坑
* etcd web ui等需要第三方 但是基本功能稳定特别是watch和读写能力 有google使用背书


### Appliction 模块
* 基于Ioc 控制反转
* 游戏逻辑单goroutine 开发简单高效
* 类似Half Sync/Half Async模式 net timer ticker单独跑goroutine --通过chan mq--> main goroutine
* 游戏逻辑靠客户端Req和服务器Tick驱动
* TimerQueue timer ticker trans to main loop

### Usage
### Module 接口
```
type Module interface {
    OnInit() bool
    OnFinal()
    Run(interval time.Duration, closeSig chan os.Signal)
}
```
####GameSvr服
![cmd-markdown-logo](game_module.png)
```

func main() {
    app := &module.Application{}
    GameModule := &GameModule{}

    app.Register(GameModule) //将GameModule 注册到app
    app.SetInterval(time.Second * DefaultInterval)
    app.Run()  //调用app.Run() 程序即进入MainLoop
}
```
