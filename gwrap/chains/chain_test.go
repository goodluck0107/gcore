package chains_test

import (
	"fmt"
	"github.com/goodluck0107/gcore/gwrap/chains"
	"testing"
)

func TestNewChain(t *testing.T) {
	c := chains.NewChain()

	defer c.FireHead()

	c.AddToHead(func() {
		fmt.Println(1111)
	})

	c.AddToHead(func() {
		fmt.Println(2222)
	})

	c.AddToHead(func() {
		fmt.Println(3333)
	})
}
