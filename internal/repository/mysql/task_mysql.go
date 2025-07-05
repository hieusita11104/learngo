package mysql

import (
	"time"

	"learngo/internal/domain"

	"gorm.io/gorm"
)

type TaskModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Done      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	CreatedBy uint      `gorm:"not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedBy uint      `gorm:"not null"`
}

func (TaskModel) TableName() string {
	return "tasks"
}

func toEntity(m *TaskModel) *domain.Task {
	return &domain.Task{
		ID:        m.ID,
		Title:     m.Title,
		Done:      m.Done,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		CreatedBy: m.CreatedBy,
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		UpdatedBy: m.UpdatedBy,
	}
}

func toModel(e *domain.Task) *TaskModel {
	return &TaskModel{
		ID:        e.ID,
		Title:     e.Title,
		Done:      e.Done,
		CreatedBy: e.CreatedBy,
		UpdatedBy: e.UpdatedBy,
	}
}

type TaskMySQLRepo struct {
	DB *gorm.DB
}

func NewTaskMySQLRepo(db *gorm.DB) domain.TaskRepository {
	if err := db.AutoMigrate(&TaskModel{}); err != nil {
		return nil
	}
	return &TaskMySQLRepo{DB: db}
}

func (r *TaskMySQLRepo) Create(t *domain.Task) error {
	model := toModel(t)
	if err := r.DB.Create(model).Error; err != nil {
		return err
	}
	t.ID = model.ID // Update the ID after creation
	return nil
}

func (r *TaskMySQLRepo) GetAll(limit, offset int, filters domain.Task) ([]domain.Task, int64, error) {
	var models []TaskModel
	query := r.DB.Model(&TaskModel{})

	if filters.Title != "" {
		query = query.Where("title LIKE ?", "%"+filters.Title+"%")
	}
	if filters.Done {
		query = query.Where("done = ?", filters.Done)
	} else if !filters.Done {
		query = query.Where("done = ?", false)
	}
	if filters.CreatedBy != 0 {
		query = query.Where("created_by = ?", filters.CreatedBy)
	}

	var total int64
	query.Model(&TaskModel{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	tasks := []domain.Task{}
	for _, m := range models {
		tasks = append(tasks, *toEntity(&m))
	}
	return tasks, total, nil
}

func (r *TaskMySQLRepo) GetByID(id uint) (*domain.Task, error) {
	var m TaskModel
	if err := r.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *TaskMySQLRepo) Update(t *domain.Task) error {
	return r.DB.Save(toModel(t)).Error
}

func (r *TaskMySQLRepo) Delete(id uint) error {
	return r.DB.Delete(&TaskModel{}, id).Error
}

// DeleteAll deletes all tasks from the database.
func (r *TaskMySQLRepo) DeleteAll() error {
	return r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&TaskModel{}).Error
}
