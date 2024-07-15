package book

// Book 定义领域类型。
type Book struct {
	ID     int
	Name   string
	Author string
}

// Reader 定义【读数据】接口。
type Reader interface {
	Get(id int) (*Book, error)
}

// Writer 定义【写数据】接口。
type Writer interface {
	Create(e *Book) (int, error)
}

// Repository 定义【持久化】接口。
type Repository interface {
	Reader
	Writer
}

// UseCase 定义【用例】接口。
type UseCase interface {
	Get(id int) (*Book, error)
	GetAll() ([]*Book, error)
	Create(title, author string, pages int, quantity int) (int, error)
}
