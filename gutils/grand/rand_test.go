package grand_test

import (
	"fmt"
	"github.com/goodluck0107/gcore/gutils/gconv"
	"github.com/goodluck0107/gcore/gutils/grand"
	"testing"
)

func Test_Str(t *testing.T) {
	t.Log(grand.Str("您好中国AJCKEKD", 5))
}

func Test_Symbols(t *testing.T) {
	t.Log(grand.Symbols(5))
}

func Test_Int(t *testing.T) {
	t.Log(grand.Int(1, 2))
}

func Test_Float32(t *testing.T) {
	t.Log(grand.Float32(-50, 5))
}

func TestLucky(t *testing.T) {
	t.Log(grand.Lucky(50.201222))
	t.Log(grand.Lucky(0.201222))
	t.Log(grand.Lucky(50))
	t.Log(grand.Lucky(0))
}

func TestWeight(t *testing.T) {
	total := 1000000

	weights := []interface{}{
		50,
		20.3,
		29.7,
	}

	counters := []int{0, 0, 0}

	for i := 0; i < total; i++ {
		index := grand.Weight(func(v interface{}) float64 {
			return gconv.Float64(v)
		}, weights...)
		counters[index] = counters[index] + 1
	}

	for i, num := range counters {
		fmt.Printf("index: %d, weight: %f, probability: %f\n", i, gconv.Float64(weights[i]), float64(num)/float64(total)*100)
	}
}
