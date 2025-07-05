package postgres

import (
	"learngo/internal/domain"

	"gorm.io/gorm"
)

type Task struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func toEntity(m *Task) *domain.Task {
	return &domain.Task{
		ID:    m.ID,
		Title: m.Title,
		Done:  m.Done,
	}
}
func toModel(e *domain.Task) *Task {
	return &Task{
		ID:    e.ID,
		Title: e.Title,
		Done:  e.Done,
	}
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *domain.Task) error {
	return r.db.Create(toModel(task)).Error
}

func (r *TaskRepository) GetByID(id uint) (*domain.Task, error) {
	var model Task
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}
func (r *TaskRepository) GetAll() ([]domain.Task, error) {
	var models []Task
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	var tasks []domain.Task
	for _, m := range models {
		tasks = append(tasks, *toEntity(&m))
	}
	return tasks, nil
}

func (r *TaskRepository) Update(task *domain.Task) error {
	model := toModel(task)
	if err := r.db.Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&Task{}, id).Error
}

func (r *TaskRepository) DeleteAll() ([]domain.Task, error) {
	var models []Task
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	if err := r.db.Delete(&models).Error; err != nil {
		return nil, err
	}
	var tasks []domain.Task
	for _, m := range models {
		tasks = append(tasks, *toEntity(&m))
	}
	return tasks, nil
}
