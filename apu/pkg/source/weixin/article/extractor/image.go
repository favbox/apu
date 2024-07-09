package extractor

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"apu/pkg/schema"
	"github.com/bytedance/gopkg/util/xxhash3"
)

var rePicturePageInfoListJs = regexp.MustCompile(`var picturePageInfoList\s*=\s*"([^"]+)"`)

type PicInfo struct {
	CdnUrl string `json:"cdn_url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

// ExtractImages 从响应体中提取图片列表。
func ExtractImages(body []byte) ([]*schema.Image, error) {
	submatch := rePicturePageInfoListJs.FindSubmatch(body)
	if len(submatch) != 2 {
		return nil, errors.New("无法提取文中图片页面信息列表，请检查文档 picturePageInfoList")
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
		return nil, err
	}

	// 遍历图片并回调函数
	var images []*schema.Image
	for _, p := range picturePageInfoList {
		width, _ := strconv.Atoi(p.Width)
		height, _ := strconv.Atoi(p.Height)
		image := &schema.Image{
			Source:      schema.Weixin,
			Width:       width,
			Height:      height,
			OriginalUrl: p.CdnUrl,
		}
		if err = packImage(image); err == nil {
			images = append(images, image)
		} else {
			log.Println(err, p.CdnUrl)
		}
	}

	return images, nil
}

func packImage(image *schema.Image) error {
	parsedURL, err := url.Parse(image.OriginalUrl)
	if err != nil {
		return err
	}
	hostname := parsedURL.Hostname()
	if hostname != "mmbiz.qpic.cn" {
		return errors.New("微信公众号图片域名异常")
	}

	imagePath := parsedURL.Path
	if !strings.HasPrefix(imagePath, "/mmbiz_") {
		return errors.New("微信公众号图片路径前缀异常")
	}
	if vs := strings.SplitN(imagePath, "/", 4); len(vs) == 4 {
		image.Format = strings.TrimPrefix(vs[1], "mmbiz_")
		if vs[2] != "" {
			image.OriginalUrl = vs[2]
			image.Key = xxhash3.HashString(vs[2])
		}
	}
	return nil
}
