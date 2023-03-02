package example

import (
	"fmt"

	"github.com/jamestack/snow"
)

// 游戏节点对象
type Game struct {
	*snow.Node[Game]
}

// 网关消息处理
func (g *Game) ReadMsg(id string, stream *snow.Stream) {
	defer stream.Close()
	// 有新链接连入，id为该链接生成的唯一id
	fmt.Println("new Receive:", id)

	go func() {
		for {
			data, err := stream.ReadData()
			if err != nil {
				fmt.Println("read data err:", err)
				return
			}

			fmt.Println("read msg:", string(data))

			// 原封不动将消息返回去，模拟回复协议
			_, err = stream.Write(data)
			if err != nil {
				fmt.Println("write data err:", err)
				return
			}
		}
	}()

}
