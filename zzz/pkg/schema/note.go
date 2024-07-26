package schema

const (
	NoteTypeNormal = "normal" // 常规图文笔记
	NoteTypeVideo  = "video"  // 视频笔记
)

const (
	NoteStateInit       = 0 // 刚采集或创建的初始笔记
	NoteStateCounted    = 1 // 交互量已采集
	NoteStateDetailed   = 2 // 详情页已采集
	NoteStateStructured = 3 // 内容已结构化
	NoteStateEmbedded   = 4 // 向量已提取
	NoteStateIndexed    = 5 // 内容已索引
)
