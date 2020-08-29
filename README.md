# snow
Snow是一个面向分布式应用的开发框架。

#### 简介
Snow为应用提供灵活易用的分布式节点抽象，
开发者无需关注节点间通信细节，仅需了解简单的概念，即可构建出功能强大的分布式应用。

#### 概念
##### 一、集群
Snow抽象了一个集群的概念，所有的Snow子节点都挂载在Snow集群之上，
集群可以由一个Master和多个Peer构成，也可以由单个Master组成。

#### 二、挂载点
Snow集群内会维护一个子节点挂载树，所有挂载在节点树上面的子节点间都可以互相通信，Snow会保证集群内Master与所有Peer的挂载点的一致性。

##### 三、子节点
子节点是Snow集群最核心的概念，也是集群内的最小业务实体，子节点间可以非常轻松的相互通信，而无需关心该节点在本地还是远程（本地节点调用无网络开销）。
##### 四、特色
1. 开发者无需关心节点间的一切通信细节，节点树是实时同步的，查询节点也无需网络开销。
2. 让开发者专注于架构设计，专注于业务开发，无需关心调用的节点在本地还是远程，节省一切与实际业务无关的代码。
3. 节点间除了支持普通rpc Call调用，还支持Stream实时流传输(几十行代码即可实现一套功能完善的分布式网管)，详情见example目录。
4. 无需第三方基础服务支持(redis、etcd等)，直接一键运行。
5. 灵活的开发模式，可以实现开发测试时用Master单服模式，部署时使用多服部署，没有做不到只有想不到。
#### 例子
A节点:
```go
package main

import (
	"fmt"
	"snow"
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
	err, done := snow.ServeMaster("127.0.0.1:8000", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("snow.ServeMaster", err)
		return
	}
	defer func() {
		err = <-done
		fmt.Println("Serve End", err)
	}()

	// 挂载到"/james"节点下
	snow.Mount("james", &Person{Name: "james"})
}
```

B节点:
```go
package main

import (
	"fmt"
	"snow"
	"time"
)

func main() {
	// 将snow设为Peer子节点启动
	err, done := snow.ServePeer("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8001")
	if err != nil {
		fmt.Println("snow.ServePeer", err)
		return
	}
	defer func() {
		err = <-done
		fmt.Println("Serve End", err)
	}()

	// 等待首次挂载点同步完
	<-time.After(3 * time.Second)

	// 查找某个要通信的节点
	james := snow.Find("/james")
	// 节点方法调用
	err = james.Call("Say", "jack", func(res string) {
		fmt.Println(res)
	})
	fmt.Println(err)
}
```
