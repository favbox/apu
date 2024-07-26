package shortid

import (
	"errors"
	"fmt"
	"strconv"

	"apu/pkg/schema"
	"apu/pkg/source/weixin"
	"apu/pkg/util/stringx"
	"github.com/sqids/sqids-go"
)

const alphabet = "0123456789abcdef"

var s *sqids.Sqids

func init() {
	s, _ = sqids.New(sqids.Options{
		Alphabet:  alphabet,
		MinLength: 24,
	})
}

func Decode(id string) (schema.Source, []uint64, error) {
	nums := s.Decode(id)
	if len(nums) < 2 {
		return 0, nil, errors.New("invalid id")
	}
	source := int(nums[0])
	switch source {
	case schema.SourceWeixin.Int():
		if len(nums) != 4 {
			return 0, nil, fmt.Errorf("invalid id for source %d", source)
		}
		return schema.SourceWeixin, nums, nil
	}

	return 0, nil, errors.New("not supported id source")
}

func EncodeWeixinArticleID(biz, mid, idx string) (string, error) {
	ghID := weixin.Biz2GhID(biz)
	u64GhID := stringx.MustNumber[uint64](fmt.Sprintf("%d", ghID))
	u64Mid := stringx.MustNumber[uint64](mid)
	u64Idx := stringx.MustNumber[uint64](idx)

	id, err := s.Encode([]uint64{uint64(schema.SourceWeixin), u64GhID, u64Mid, u64Idx})
	if err != nil {
		return "", err
	}
	return id, nil
}

func DecodeWeixinArticleID(id string) (biz, mid, idx string, err error) {
	nums := s.Decode(id)
	if len(nums) != 4 {
		return "", "", "", errors.New("invalid id")
	}
	if nums[0] != uint64(schema.SourceWeixin) {
		return "", "", "", errors.New("invalid source")
	}

	i64GhID := int64(nums[1])
	biz = weixin.GhID2Biz(i64GhID)
	mid = strconv.FormatUint(nums[2], 10)
	idx = strconv.FormatUint(nums[3], 10)

	return biz, mid, idx, nil
}
