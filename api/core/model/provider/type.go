package provider

import "apu/core/model"

type (
	// ConfigureMethodEnum 配置方法枚举。
	ConfigureMethodEnum string
	configureMethod     struct {
		PredefinedModel   ConfigureMethodEnum
		CustomizableModel ConfigureMethodEnum
	}
)

// ConfigureMethod 供应商模型的配置方法。
var ConfigureMethod = configureMethod{
	PredefinedModel:   "predefined-model",
	CustomizableModel: "customizable-model",
}

type (
	// FormTypeEnum 表单类型枚举。
	FormTypeEnum string
	formType     struct {
		TextInput   FormTypeEnum
		SecretInput FormTypeEnum
		Select      FormTypeEnum
		Radio       FormTypeEnum
		Switch      FormTypeEnum
	}
)

// FormType 表单类型。
var FormType = formType{
	TextInput:   "text-input",
	SecretInput: "secret-input",
	Select:      "select",
	Radio:       "radio",
	Switch:      "switch",
}

// FormShowOnObject 用于表单显示。
type FormShowOnObject struct {
	Variable string
	Value    string
}

// FormOption 表单选项。
type FormOption struct {
	Label  *model.I18n
	Value  string
	ShowOn []FormShowOnObject
}

func NewFormOption(value string, showOn []FormShowOnObject, label *model.I18n) *FormOption {
	opt := &FormOption{
		Label:  label,
		Value:  value,
		ShowOn: showOn,
	}
	if opt.Label == nil {
		opt.Label = &model.I18n{
			EnUS: opt.Value,
		}
	}
	return opt
}
