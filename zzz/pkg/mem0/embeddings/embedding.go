package embeddings

import "context"

type Embedder interface {
	// Embed 获取单个文本的嵌入。
	Embed(ctx context.Context, text string) ([]float32, error)
}
