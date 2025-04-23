package fastjson

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestOom(t *testing.T) {
	f, err := os.Open("./oom_test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	b, err := io.ReadAll(reader)
	// ParseBytes 方法调用后，内存急剧上升，大概500-600M
	GetStringImproved(b, []string{"data", "token"}...)

	time.Sleep(time.Minute)
}
