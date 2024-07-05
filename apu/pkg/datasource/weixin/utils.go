package weixin

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

func Biz2BookId(biz string) string {
	return GhId2BookId(Biz2GhId(biz))
}

func Biz2GhId(biz string) int64 {
	ghIdBs, err := base64.StdEncoding.DecodeString(biz)
	if err != nil {
		return 0
	}
	ghId, err := strconv.ParseInt(string(ghIdBs), 10, 64)
	if err != nil {
		return 0
	}
	return ghId
}

func GhId2BookId(ghId int64) string {
	return fmt.Sprintf("MP_WXS_%d", ghId)
}

func GhId2Biz(ghId int64) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", ghId)))
}

func BookId2GhId(bookId string) int64 {
	if !strings.HasPrefix(bookId, "MP_WXS_") {
		return 0
	}
	ghId, err := strconv.ParseInt(bookId[7:], 10, 64)
	if err != nil {
		return 0
	}
	return ghId
}
