package coze_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"apu/pkg/llm/coze"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func TestStructureProject(t *testing.T) {
	//	query := `项目名称 | 保利海晏天珺
	//项目地址 | 中国 · 宁波
	//完成时间 | 2024年
	//设计面积｜185㎡、155 ㎡、145 ㎡
	// 保利发展华东区域共享中心
	//
	//               浙江保利城市发展有限公司
	//
	//建筑设计 | 广州博瀚建筑设计事务所有限公司
	//景观设计 | 加特林（重庆）景观规划设计有限公司
	//室内设计 | MAUDEA 牧笛设计
	//软装设计 | MAUDEA 牧笛设计
	//项目摄影 | 鲁哈哈`
	result, err := coze.StructureProject(query)
	assert.Nil(t, err)
	for _, msg := range result.Messages {
		if msg.Type != "answer" {
			continue
		}
		var content coze.Content
		err := json.Unmarshal([]byte(msg.Content), &content)
		assert.Nil(t, err)

		var output coze.Output
		err = json.Unmarshal([]byte(content.Output), &output)
		assert.Nil(t, err)

		fmt.Println("时间", output.Time)
		fmt.Println("名称", output.Name)
		fmt.Println("标题", output.Title)
		fmt.Println("关键词", output.Keyword)
		fmt.Println("面积", output.Area)
		fmt.Println("地点", output.Position)
		fmt.Println("分类", output.Category)
		fmt.Println("设计", output.Designer)
	}
}

//go:embed prompt.tpl
var prompt string

//go:embed query.tpl
var query string

func TestOllamaAndQwen2(t *testing.T) {
	ctx := context.Background()
	llm, err := ollama.New(
		ollama.WithModel("qwen2"),
		ollama.WithFormat("json"),
	)
	if err != nil {
		log.Fatal(err)
	}

	prompt := strings.ReplaceAll(prompt, "{{input}}", query)

	//content := []llms.MessageContent{
	//	llms.TextParts(llms.ChatMessageTypeSystem, prompt),
	//	llms.TextParts(llms.ChatMessageTypeHuman, query),
	//}

	////completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	////	fmt.Println(string(chunk))
	////	return nil
	////}))
	//completion, err := llm.GenerateContent(ctx, content)

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
	//for _, c := range completion.Choices {
	//	fmt.Println(c.Content)
	//}
	//_ = completion
}
