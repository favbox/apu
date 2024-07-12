package extractor

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"

	"apu/pkg/schema"
	"apu/pkg/source"
	"apu/pkg/util/stringx"
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
		return ExtractImagesV2(body)
		//return nil, nil, errors.New("无法提取文中图片页面信息列表，请检查文档 picturePageInfoList")
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

// ExtractImagesV2 公众号仿小红书的图文版本
func ExtractImagesV2(body []byte) ([]*schema.Image, map[uint64][2]int, error) {
	body = bytes.ReplaceAll(body, []byte("\\x26amp;"), []byte("&"))
	body = bytes.ReplaceAll(body, []byte("&amp;"), []byte("&"))

	start := bytes.Index(body, []byte("window.picture_page_info_list ="))
	end := bytes.Index(body[start:], []byte("slice(0, 20);"))
	picturePageInfoList := string(body[start : start+end])

	var cdnUrls []string
	re := regexp.MustCompile(`cdn_url:\s*'(.*?)',`)
	matches := re.FindAllStringSubmatch(picturePageInfoList, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			cdnUrls = append(cdnUrls, match[1])
		}
	}

	var widths []int
	re = regexp.MustCompile(`width:\s*'(\d+)'`)
	matches = re.FindAllStringSubmatch(picturePageInfoList, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			widths = append(widths, stringx.MustNumber[int](match[1]))
		}
	}
	var heights []int
	re = regexp.MustCompile(`height:\s*'(\d+)'`)
	matches = re.FindAllStringSubmatch(picturePageInfoList, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			heights = append(heights, stringx.MustNumber[int](match[1]))
		}
	}

	var images []*schema.Image
	var imageSizeMap map[uint64][2]int
	if len(cdnUrls) == len(widths) && len(widths) == len(heights) {
		imageSizeMap = make(map[uint64][2]int, len(cdnUrls))
		for i := 0; i < len(cdnUrls); i++ {
			image := &schema.Image{
				Source:      schema.Weixin,
				Key:         source.Key(cdnUrls[i]),
				OriginalUrl: cdnUrls[i],
				Width:       widths[i],
				Height:      heights[i],
			}
			images = append(images, image)
			imageSizeMap[image.Key] = [2]int{image.Width, image.Height}
		}
	}

	return images, imageSizeMap, nil
}
