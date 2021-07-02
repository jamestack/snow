# snow
Snow是一个面向分布式应用的开发框架。

#### 简介
Snow为应用提供灵活易用的分布式节点抽象，
开发者无需关注节点间通信细节，仅需了解简单的概念，即可构建出功能强大的分布式应用。

#### 概念
##### 一、集群
Snow抽象了一个集群的概念，所有的Snow子节点都挂载在Snow集群之上，同一集群内的节点之间通信无网络开销。

##### 二、服务与子节点
Snow集群内可以挂载任意多个服务，每个服务下可能会有一个或者多个子节点：
```shell script
# 子节点命名规则：
ServerName/NodeName    # ServerName/NodeName定位唯一子节点
ServerName             # 如果该服务只有一个子节点，那么可以缩写成ServerName的方式，Snow会自动将其补全为ServerName/Master
```

##### 三、子节点
1. 子节点是Snow集群最核心的概念，也是集群内的最小业务实体。
2. 子节点间可以非常轻松的相互通信，而无需关心该节点在本地还是远程（本地节点调用无网络开销）。
3. 子节点间除了支持普通rpc Call调用，还支持Stream实时流传输(几十行代码即可实现一套功能完善的分布式网关)，详情见example目录。
##### 四、特色
1. 开发者无需关心节点间的一切通信细节，也无需关心Snow与Consul、Etcd等第三方发现服务之间的交互细节。
2. 让开发者专注于架构设计，专注于业务开发，无需关心调用的节点在本地还是远程，节省一切与实际业务无关的代码。
3. 集群支持免Consul单机单集群运行方便测试，也可以接入Consul单机多集群或者多机多集群部署。
4. Snow鼓励开发者将子节点模块化、组件化，比如抽象出Auth模块、网关模块、游戏服务模块、广播模块等，进一步压缩重用非业务代码。
#### 例子
A节点:
```go
package main

import (
	"fmt"
	"github.com/jamestack/snow"
)

// 每个Struct只需包含*snow.Node，那么Snow就认为这是一个子节点
type Person struct {
	Name string
	*snow.Node
}

// 定义一个Say方法
func (p *Person) Say(name string) string {
	return "hello " + name + ", my name is " + name + "."
}

func main() {
	cluster := snow.NewClusterWithConsul("127.0.0.1:8000", "127.0.0.1:8000")
	done, err := cluster.Serve()
	if err != nil {
		fmt.Println("snow.ServeMaster", err)
		return
	}
	defer func() {
		<-done
		fmt.Println("Serve End", err)
	}()

	// 挂载到"/james"节点下
	cluster.Mount("james", &Person{Name: "james"})
}
```

B节点:
```go
package main

import (
	"fmt"
	"github.com/jamestack/snow"
	"time"
)

func main() {
	cluster := snow.NewClusterWithConsul("127.0.0.1:8001", "127.0.0.1:8001")
	done, err := cluster.Serve()
	if err != nil {
		fmt.Println("snow.ServeMaster", err)
		return
	}
	defer func() {
		<-done
		fmt.Println("Serve End", err)
	}()

	// 等待首次挂载点同步完
	<-time.After(3 * time.Second)

	// 查找某个要通信的节点
	james, err := cluster.Find("james")
	fmt.Println(err)
	// 节点方法调用
	err = james.Call("Say", "jack", func(res string) {
		fmt.Println(res)
	})
	fmt.Println(err)
}
```
