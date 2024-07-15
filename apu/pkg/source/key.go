package source

import "github.com/bytedance/gopkg/util/xxhash3"

// UniqueID 返回数据源下的唯一编号。
func UniqueID(s string) uint64 {
	return xxhash3.HashString(s)
}
