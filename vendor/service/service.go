package service

import  (
    "net"
    "github.com/coreos/etcd/client"
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "log"
)

const (
    DefaultServicePath = "/services"
    DefualtServiceName = "services/names"
    DEFAULT_ETCD_HOST  = "http://172.17.0.2:2379"
)

type node struct {
    key  string
    conn *net.TCPConn
}
type service struct {
    nodes   []*node
    idx     uint32
}
type ServiceMgr struct {
    services        map[string]*service
    known_names     map[string]bool
    addServiceMQ    chan node
    removeServiceMQ chan string
    etcdClient      client.Client
}

func (sm *ServiceMgr) OnInit() {

    machines := []string{DEFAULT_ETCD_HOST}
    if env := os.Getenv("ETCD_HOST"); env != "" {
        machines = strings.Split(env, ";")
    }

    // init etcd client
    cfg := client.Config{
        Endpoints: machines,
        Transport: client.DefaultTransport,
    }

    cli, err := client.New(cfg)
    if err != nil {
        fmt.Println(err)
    }
    sm.etcdClient = cli

    sm.services = make(map[string]*service)
    sm.known_names = make(map[string]bool)
    sm.addServiceMQ = make(chan node, 10)
    sm.removeServiceMQ = make(chan string, 10)

    names := sm.loadNames()

    for _, name := range names {
        sm.known_names[DefaultServicePath + "/" + strings.TrimSpace(name)] = true
    }

    sm.connectAll(DefaultServicePath)
}

func (sm *ServiceMgr) Select() {
    select {
    case node := <-sm.addServiceMQ:
        sm.addService(node.key, node.conn)
    case key := <-sm.removeServiceMQ:
        sm.removeService(key)
    default:
    }
}

func (sm *ServiceMgr) watcher() {
    kApi := client.NewKeysAPI(sm.etcdClient)
    watcher := kApi.Watcher(DefaultServicePath, &client.WatcherOptions{Recursive:true})

    for {
        rsp, err := watcher.Next(context.Background())
        if err != nil {
            fmt.Println(err)
            continue
        }
        fmt.Println(rsp.Action)
        if rsp.Node.Dir {
            fmt.Println("not file")
            continue
        }

        switch rsp.Action {
        case "set", "create", "update", "compareAndSwap":
            tcpAddr, err := net.ResolveTCPAddr("tcp4", rsp.Node.Value)
            if err != nil {
                fmt.Println(err)
            }
            conn, err := net.DialTCP("tcp", nil, tcpAddr)
            if err == nil {
                sm.addServiceMQ <- node{rsp.Node.Key, conn}
            } else {
                log.Println(err)
                fmt.Println("can not connect ",rsp.Node.Key, rsp.Node.Value)
            }

        case "delete":
            sm.removeServiceMQ <- rsp.Node.Key
        
        }
    }
}

func (sm *ServiceMgr) loadNames() []string {
    kApi := client.NewKeysAPI(sm.etcdClient)
    fmt.Println("reading names :", DefualtServiceName)
    rsp, err := kApi.Get(context.Background(), DefualtServiceName, nil)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    if rsp.Node.Dir {
        fmt.Println("names is not a file")
    }

    return strings.Split(rsp.Node.Value, "/n")
}

func (sm *ServiceMgr) connectAll(dir string) {
    kApi := client.NewKeysAPI(sm.etcdClient)
    fmt.Println("connecting services under:", dir)
    rsp, err := kApi.Get(context.Background(), dir, &client.GetOptions{Recursive: true})
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, node := range rsp.Node.Nodes {
        if node.Dir {
            for _, service := range node.Nodes {
                service_name := filepath.Dir(service.Key)
                if !sm.known_names[service_name] {
                    continue
                }
                tcpAddr, err := net.ResolveTCPAddr("tcp4", rsp.Node.Value)
                if err != nil {
                    fmt.Println(err)
                    continue
                }
                conn, err := net.DialTCP("tcp", nil, tcpAddr)
                if err == nil {
                    sm.addService(service.Key, conn)
                    fmt.Println("connect service:" + service.Value)
                }
            }
        }
    }
    fmt.Println("services add complete")
    go sm.watcher()
}

func (sm *ServiceMgr) addService(key string, conn *net.TCPConn) {
    fmt.Println("add service" + key)
    serviceName := filepath.Dir(key)
    if !sm.known_names[serviceName] {
        return
    }

    if sm.services[serviceName] == nil {
        sm.services[serviceName] = &service{}
    }

    service := sm.services[serviceName]
    node := &node{
        key: key,
        conn: conn,
    }
    service.nodes = append(service.nodes, node)
}

func (sm *ServiceMgr) removeService(key string) {
    serviceName := filepath.Dir(key)
    if !sm.known_names[serviceName] {
        return
    }
    service := sm.services[serviceName]
    if service == nil {
        fmt.Println("no such service:", serviceName)
        return
    }

    for idx := range service.nodes {
        if service.nodes[idx].key == key {
            service.nodes[idx].conn.Close()
            service.nodes = append(service.nodes[:idx], service.nodes[idx+1:]...)
            fmt.Println("service remove:", key)
            return
        }
    }
}



