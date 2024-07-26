package memory

// Memorizer 定义记忆体需要实现的方法。
type Memorizer interface {
	Get(memoryID string)
	GetAll()
	Update(memoryID string, date any)
	Delete(memoryID string)
	History(memoryID string)
}
