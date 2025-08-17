package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"todolist/internal/model"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
)

// TodoRepository 待办事项仓储接口
type TodoRepository interface {
	GetAll() ([]*model.Todo, error)
	GetByID(id string) (*model.Todo, error)
	Create(todo *model.Todo) error
	Update(todo *model.Todo) error
	Delete(id string) error
	GetByCompleted(completed bool) ([]*model.Todo, error)
}

// MemoryTodoRepository 内存实现的待办事项仓储
type MemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[string]*model.Todo
}

// NewMemoryTodoRepository 创建新的内存待办事项仓储
func NewMemoryTodoRepository() TodoRepository {
	repo := &MemoryTodoRepository{
		todos: make(map[string]*model.Todo),
	}
	
	// 初始化示例数据
	repo.initSampleData()
	
	return repo
}

func (r *MemoryTodoRepository) initSampleData() {
	now := time.Now()
	
	todos := []*model.Todo{
		{
			ID:          uuid.New().String(),
			Title:       "学习Go语言",
			Description: "完成Go语言基础教程",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New().String(),
			Title:       "重构代码",
			Description: "按照Go项目标准结构重构todolist代码",
			Completed:   true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	
	for _, todo := range todos {
		r.todos[todo.ID] = todo
	}
}

func (r *MemoryTodoRepository) GetAll() ([]*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	todos := make([]*model.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	
	return todos, nil
}

func (r *MemoryTodoRepository) GetByID(id string) (*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}
	
	return todo, nil
}

func (r *MemoryTodoRepository) Create(todo *model.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	todo.ID = uuid.New().String()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	
	r.todos[todo.ID] = todo
	return nil
}

func (r *MemoryTodoRepository) Update(todo *model.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.todos[todo.ID]; !exists {
		return ErrTodoNotFound
	}
	
	todo.UpdatedAt = time.Now()
	r.todos[todo.ID] = todo
	return nil
}

func (r *MemoryTodoRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.todos[id]; !exists {
		return ErrTodoNotFound
	}
	
	delete(r.todos, id)
	return nil
}

func (r *MemoryTodoRepository) GetByCompleted(completed bool) ([]*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var todos []*model.Todo
	for _, todo := range r.todos {
		if todo.Completed == completed {
			todos = append(todos, todo)
		}
	}
	
	return todos, nil
}