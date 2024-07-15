package weixin

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

func Biz2BookID(biz string) string {
	return GhId2BookID(Biz2GhID(biz))
}

func Biz2GhID(biz string) int64 {
	ghIDBytes, err := base64.StdEncoding.DecodeString(biz)
	if err != nil {
		return 0
	}
	ghID, err := strconv.ParseInt(string(ghIDBytes), 10, 64)
	if err != nil {
		return 0
	}
	return ghID
}

func GhId2BookID(ghID int64) string {
	return fmt.Sprintf("MP_WXS_%d", ghID)
}

func GhID2Biz(ghID int64) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", ghID)))
}

func BookID2GhID(bookID string) int64 {
	if !strings.HasPrefix(bookID, "MP_WXS_") {
		return 0
	}
	ghID, err := strconv.ParseInt(bookID[7:], 10, 64)
	if err != nil {
		return 0
	}
	return ghID
}
