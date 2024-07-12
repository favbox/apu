package extractor

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"

	"apu/pkg/schema"
	"apu/pkg/source"
)

var rePicturePageInfoListJs = regexp.MustCompile(`var picturePageInfoList\s*=\s*"([^"]+)"`)

type PicInfo struct {
	CdnUrl string `json:"cdn_url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

// ExtractImages 从响应体中提取图片列表。
func ExtractImages(body []byte) ([]*schema.Image, map[uint64][2]int, error) {
	submatch := rePicturePageInfoListJs.FindSubmatch(body)
	if len(submatch) != 2 {
		return nil, nil, errors.New("无法提取文中图片页面信息列表，请检查文档 picturePageInfoList")
	}

	jsonBytes := submatch[1]

	// 规范化 json 字符串
	if bytes.Contains(jsonBytes, []byte(",]")) {
		jsonBytes = bytes.ReplaceAll(jsonBytes, []byte((",]")), []byte("]"))
	}
	jsonBytes = bytes.ReplaceAll(jsonBytes, []byte("'"), []byte("\""))
	jsonBytes = bytes.ReplaceAll(jsonBytes, []byte("\\x26amp;"), []byte("&"))
	jsonBytes = bytes.ReplaceAll(jsonBytes, []byte("&amp;"), []byte("&"))

	// 反序列化图片列表
	var picturePageInfoList []PicInfo
	err := json.Unmarshal(jsonBytes, &picturePageInfoList)
	if err != nil {
		return nil, nil, err
	}

	// 追加图片及唯一键
	var images []*schema.Image
	imageSizeMap := make(map[uint64][2]int, len(picturePageInfoList))
	for _, p := range picturePageInfoList {
		width, _ := strconv.Atoi(p.Width)
		height, _ := strconv.Atoi(p.Height)
		image := &schema.Image{
			Source:      schema.Weixin,
			Width:       width,
			Height:      height,
			OriginalUrl: p.CdnUrl,
			Key:         source.Key(p.CdnUrl),
		}
		images = append(images, image)
		imageSizeMap[image.Key] = [2]int{image.Width, image.Height}
	}

	return images, imageSizeMap, nil
}
