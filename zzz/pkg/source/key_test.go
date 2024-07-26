package source

import (
	"fmt"
	"testing"
)

func TestKey(t *testing.T) {
	id := UniqueID("https://sqids.org/zh")
	fmt.Println(id)
}
