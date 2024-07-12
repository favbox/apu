package source

import "github.com/bytedance/gopkg/util/xxhash3"

// Key 返回数据源下的唯一键。
func Key(s string) uint64 {
	return xxhash3.HashString(s)
}
