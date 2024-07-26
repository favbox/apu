package model

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Type 模型类型。
type Type uint8

const (
	TypeLLM Type = iota + 1
	TypeTextEmbedding
	TypeRerank
	TypeSpeech2Text
	TypeModeration
	TypeTTS
	TypeText2Img
)

// 获取 model.Type 的原始模型类型。
func (t Type) String() string {
	return [...]string{"text-generation", "embeddings", "reranking", "speech2text", "tts", "moderation", "text2img"}[t-1]
}

// TypeFrom 获取原始模型类型的 model.Type。
func TypeFrom(s string) (Type, error) {
	switch s {
	case "text-generation", "llm":
		return TypeLLM, nil
	case "embeddings":
		return TypeTextEmbedding, nil
	case "reranking", "rerank":
		return TypeRerank, nil
	case "speech2text":
		return TypeSpeech2Text, nil
	case "tts":
		return TypeTTS, nil
	case "text2img":
		return TypeText2Img, nil
	default:
		return 0, fmt.Errorf("无效的原始模型类型 %s", s)
	}
}

type (
	// FetchFromEnum 模型的获取来源类型。
	FetchFromEnum string
	fetchFrom     struct {
		// 来自预训练模型
		PredefinedModel FetchFromEnum
		// 来自自定义模型
		CustomizableModel FetchFromEnum
	}
)

// FetchFrom 模型的获取来源。
var FetchFrom = fetchFrom{
	PredefinedModel:   "predefined-model",
	CustomizableModel: "customizable-model",
}

type (
	FeatureEnum string
	feature     struct {
		ToolCall       FeatureEnum
		MultiToolCall  FeatureEnum
		StreamToolCall FeatureEnum
		AgentThought   FeatureEnum
		Vision         FeatureEnum
	}
)

// Feature 定义LLM大语言模型的特性。
var Feature = feature{
	ToolCall:       "tool-call",
	MultiToolCall:  "multi-tool-call",
	StreamToolCall: "stream-tool-call",
	AgentThought:   "agent-thought",
	Vision:         "vision",
}

type (
	// DefaultParameterNameEnum 模型的默认参数名称类型。
	DefaultParameterNameEnum string
	defaultParameterName     struct {
		Temperature      DefaultParameterNameEnum // 稳定性
		TopP             DefaultParameterNameEnum
		PresencePenalty  DefaultParameterNameEnum // 存在惩罚
		FrequencyPenalty DefaultParameterNameEnum // 频率惩罚
		MaxTokens        DefaultParameterNameEnum // 最大令牌数
		ResponseFormat   DefaultParameterNameEnum // 响应格式
	}
)

// DefaultParameterName 模型的默认参数名称。
var DefaultParameterName = defaultParameterName{
	Temperature:      "temperature",
	TopP:             "top_p",
	PresencePenalty:  "presence_penalty",
	FrequencyPenalty: "frequency_penalty",
	MaxTokens:        "max_tokens",
	ResponseFormat:   "response_format",
}

// DefaultParameterNameFrom 从字符串获取参数名称。
func DefaultParameterNameFrom(s string) (DefaultParameterNameEnum, error) {
	switch s {
	case "temperature":
		return DefaultParameterName.Temperature, nil
	case "top_p":
		return DefaultParameterName.TopP, nil
	case "presence_penalty":
		return DefaultParameterName.PresencePenalty, nil
	case "frequency_penalty":
		return DefaultParameterName.FrequencyPenalty, nil
	case "max_tokens":
		return DefaultParameterName.MaxTokens, nil
	case "response_format":
		return DefaultParameterName.ResponseFormat, nil
	default:
		return "", fmt.Errorf("无效的参数名称 %s", s)
	}
}

type (
	// ParameterTypeEnum 模型的参数类型类型。
	ParameterTypeEnum string
	parameterType     struct {
		Float   ParameterTypeEnum
		Int     ParameterTypeEnum
		String  ParameterTypeEnum
		Boolean ParameterTypeEnum
	}
)

// ParameterType 模型的参数类型。
var ParameterType = parameterType{
	Float:   "float",
	Int:     "int",
	String:  "string",
	Boolean: "boolean",
}

type (
	// PropertyKeyEnum 模型的属性键类型。
	PropertyKeyEnum string
	propertyKey     struct {
		Mode                    PropertyKeyEnum
		ContextSize             PropertyKeyEnum
		MaxChunks               PropertyKeyEnum
		FileUploadLimit         PropertyKeyEnum
		SupportedFileExtensions PropertyKeyEnum
		MaxCharactersPerChunk   PropertyKeyEnum
		DefaultVoice            PropertyKeyEnum
		Voices                  PropertyKeyEnum
		WordLimit               PropertyKeyEnum
		AudioType               PropertyKeyEnum
		MaxWords                PropertyKeyEnum
	}
)

// PropertyKey 模型的属性键名。
var PropertyKey = propertyKey{
	Mode:                    "mode",
	ContextSize:             "context_size",
	MaxChunks:               "max_chunks",
	FileUploadLimit:         "file_upload_limit",
	SupportedFileExtensions: "support_file_extensions",
	MaxCharactersPerChunk:   "max_characters_per_chunk",
	DefaultVoice:            "default_voice",
	Voices:                  "voices",
	WordLimit:               "word_limit",
	AudioType:               "audio_type",
	MaxWords:                "max_workers",
}

type I18n struct {
	EnUS string `json:"en_us"`
	ZhCN string `json:"zh_cn"`
}

// ProviderModel 模型供应商提供的模型。
type ProviderModel struct {
	Model           string
	Label           I18n
	ModelType       Type
	Features        *[]FeatureEnum
	FetchFrom       FetchFromEnum
	ModelProperties map[PropertyKeyEnum]any
	Deprecated      bool
	ModelConfig     map[string]any
}

// ParameterRule 模型的参数规则。
type ParameterRule struct {
	Name        string
	UseTemplate *string
	Label       I18n
	Type        ParameterTypeEnum
	Help        *I18n
	Required    bool
	Default     any
	Min         *float32
	Max         *float32
	Precision   *int
	Options     []string
}

// PriceConfig 模型的价格信息。
type PriceConfig struct {
	Input    decimal.Decimal
	Output   *decimal.NullDecimal
	Unit     decimal.Decimal
	Currency string
}

// AIModel 供应商的AI模型。
type AIModel struct {
	ProviderModel
	ParameterRules []ParameterRule
	Pricing        *PriceConfig
}

type Usage struct{}

type (
	// PriceTypeEnum 模型的价格类型
	PriceTypeEnum string
	priceType     struct {
		Input  PriceTypeEnum
		Output PriceTypeEnum
	}
)

// PriceType 模型的价格类型
var PriceType = priceType{
	Input:  "input",
	Output: "output",
}

// PriceInfo 模型的价格信息。
type PriceInfo struct {
	UnitPrice   decimal.Decimal
	Unit        decimal.Decimal
	TotalAmount decimal.Decimal
	Currency    string
}
