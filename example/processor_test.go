package example

import (
	"fmt"
	"snow"
	"testing"
	"time"
)

func TestProcessor(t *testing.T) {
	processor := snow.NewProcessor(10)

	processor.RunAfter(2*time.Second, func() {
		fmt.Println(456)
	})
	processor.Run(func() {
		fmt.Println(123)
	})

	<-time.After(3*time.Second)
}
