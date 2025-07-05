package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"learngo/internal/domain"
	"learngo/pkg/cache"
)

type TaskUsecase struct {
	Repo domain.TaskRepository
}

func NewTaskUsecase(r domain.TaskRepository) *TaskUsecase {
	return &TaskUsecase{Repo: r}
}

func (uc *TaskUsecase) CreateTask(t *domain.Task) error {
	return uc.Repo.Create(t)
}

func (uc *TaskUsecase) GetAll(limit, offset int, filter domain.Task) ([]domain.Task, int64, error) {
	return uc.Repo.GetAll(limit, offset, filter)
}

func (uc *TaskUsecase) GetTaskByID(id uint) (*domain.Task, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("task:%d", id)

	cachedTask, err := cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var task domain.Task
		if err := json.Unmarshal([]byte(cachedTask), &task); err == nil {
			return &task, nil
		}
	}
	task, err := uc.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	taskJSON, _ := json.Marshal(task)
	cache.RedisClient.Set(ctx, cacheKey, taskJSON, 10*time.Minute)

	return task, nil
}

func (uc *TaskUsecase) UpdateTask(t *domain.Task) error {
	return uc.Repo.Update(t)
}

// func (uc *TaskUsecase) DeleteTask(id uint) error {
// 	return uc.Repo.Delete(id)
// }

func (uc *TaskUsecase) DeleteAll() error {
	return uc.Repo.DeleteAll()
}
