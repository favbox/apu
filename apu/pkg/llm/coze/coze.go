package coze

import (
	"errors"

	"github.com/imroc/req/v3"
)

const accessToken = "pat_sbH8HXdFVgbEpaOC6Kcrt31MaTR1u85zdbrBUmeH8lbRg1Fh2OL1CYGxakQGGSRb"

var api *req.Client

func init() {
	api = req.C()
	api.SetCommonBearerAuthToken(accessToken)
	api.SetCommonHeader("Content-Type", "application/json")
	api.SetCommonHeader("Accept", "*/*")
	api.SetCommonHeader("Host", "api.coze.cn")
	api.SetCommonHeader("Connection", "keep-alive")
}

func StructureProject(str string) (*Result, error) {
	var result Result
	resp, err := api.R().
		SetBodyJsonMarshal(map[string]any{
			"conversation_id": "123",
			"bot_id":          "7389447004877078579",
			"user":            "123456789",
			"query":           str,
			"stream":          false,
		}).
		SetSuccessResult(&result).
		Post("https://api.coze.cn/open_api/v2/chat")
	if err != nil {
		return nil, err
	}
	if resp.IsErrorState() {
		return nil, errors.New(resp.GetStatus())
	}
	for _, m := range result.Messages {
		if m.Type != "answer" {
			continue
		}

	}

	return &result, nil
}

type Message struct {
	Role        string `json:"role"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

type Content struct {
	Output string `json:"output"`
}

type Output struct {
	Area     string `json:"area"`
	Category string `json:"category"`
	Designer string `json:"designer"`
	Keyword  string `json:"keyword"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Title    string `json:"title"`
	Time     string `json:"time"`
}

type Result struct {
	Messages       []*Message `json:"messages"`
	ConversationId string     `json:"conversation_id"`
	Code           int        `json:"code"`
	Msg            string     `json:"msg"`
}
