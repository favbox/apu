package book

var _ UseCase = (*Service)(nil)

// Service 实现 UseCase
type Service struct {
	repo Repository
}

// NewService 创建新服务。
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// Create 创建一本书。
func (s *Service) Create(title, author string, pages int, quantity int) (int, error) {
	b := &Book{
		Name:   title,
		Author: author,
	}

	if _, err := s.repo.Create(b); err != nil {
		return 0, err
	}

	return 1, nil
}

// Get 获取一本书。
func (s *Service) Get(id int) (*Book, error) {
	b, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetAll 获取一组书。
func (s *Service) GetAll() ([]*Book, error) {
	bb := []*Book{
		{
			ID:     1,
			Name:   "Book 1",
			Author: "Author 1",
		},
		{
			ID:     2,
			Name:   "Book 2",
			Author: "Author 2",
		},
	}

	return bb, nil
}
