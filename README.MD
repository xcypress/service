##服务框架
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
    serviceMgr.Init()
    fmt.Println("servicemgr start")
    for {
        serviceMgr.OnTick()
    }
}
```