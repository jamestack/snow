package example

import (
	"fmt"
	"testing"

	"github.com/jamestack/snow"
)

func TestClusterMountMaster(t *testing.T) {
	master := snow.NewClusterMaster("127.0.0.1:8000", "127.0.0.1:8000")

	done, err := master.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	master.Mount("test", &User{
		name: "james",
	})

	<-done
}

func TestClusterMountSlave(t *testing.T) {
	master := snow.NewClusterSlave("127.0.0.1:8001", "127.0.0.1:8001", "127.0.0.1:8000")

	done, err := master.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	node, err := master.Find("test")
	if err != nil {
		fmt.Println(err)
	}

	node.CallAsync("Name").Then(func(res string) {
		fmt.Println("name:", res)
	})

	<-done
}
