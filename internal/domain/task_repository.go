package domain

type TaskRepository interface {
	Create(task *Task) error
	GetByID(id uint) (*Task, error)
	GetAll(limit, offset int, filter Task) ([]Task, int64, error)
	Update(task *Task) error
	Delete(id uint) error
	DeleteAll() error
}
