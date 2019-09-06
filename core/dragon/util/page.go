package util

import (
	"strconv"
)

const (
	// todo 设置最大数据条数，目前暂定100w
	MaxPageSize = 1000000 //最大的查询页数
)

//初始化page和pageSize, 默认page为0, pageSize=10, 返回默认为字符串，方便之后查询拼接字符串
func InitPageAndPageSize(page string, pageSize string) (offset int, size int) {
	p, err := strconv.ParseInt(page, 10, 32)
	ps, errz := strconv.ParseInt(pageSize, 10, 32)
	if err != nil || errz != nil {
		return 0, 10
	}
	if p < 0 || ps < 0 || ps > MaxPageSize {
		return 0, 10
	}
	offset = int(p * ps)
	size = int(ps)
	return
}
