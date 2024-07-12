package schema

// NoteType 定义了笔记的类型
type NoteType string

const (
	NoteTypeNormal NoteType = "normal" // 常规图文笔记
	NoteTypeVideo  NoteType = "video"  // 视频笔记
)

// NoteState 定义了笔记的状态
type NoteState int32

const (
	NoteStateInit            NoteState = 0 // 刚采集或创建的初始笔记
	NoteStateInteractCrawled NoteState = 1 // 交互量已采集
	NoteStateDetailCrawled   NoteState = 2 // 详情页已采集
	NoteStateSummarized      NoteState = 3 // 摘要已经提取
	NoteStateEmbedded        NoteState = 4 // 向量已提取
)
