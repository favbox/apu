package model_test

import (
	"testing"

	"apu/core/model"
	"github.com/stretchr/testify/assert"
)

func TestModelType(t *testing.T) {
	assert.Equal(t, "text-generation", model.TypeLLM.String())
	modelType, err := model.TypeFrom("text-generation")
	assert.Nil(t, err)
	assert.Equal(t, model.TypeLLM, modelType)

	assert.Equal(t, "predefined-model", string(model.FetchFrom.PredefinedModel))
	assert.Equal(t, "tool-call", string(model.Feature.ToolCall))
	assert.Equal(t, "temperature", string(model.DefaultParameterName.Temperature))
}
