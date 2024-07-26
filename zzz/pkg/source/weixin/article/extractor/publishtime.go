package extractor

import (
	"regexp"
	"strconv"
	"time"

	"github.com/yuin/goldmark/util"
)

// 时间提取规则
var reCreateTime = regexp.MustCompile(`var oriCreateTime = '(\d+)';`)

func ExtractPublishTime(body []byte) (time.Time, bool) {
	if v := reCreateTime.FindSubmatch(body); len(v) == 2 {
		timestamp := util.BytesToReadOnlyString(v[1])
		if sec, err := strconv.ParseFloat(timestamp, 64); err == nil {
			return time.Unix(int64(sec), 0), true
		}
	}
	return time.Unix(0, 0), false
}
