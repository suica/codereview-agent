package service

import (
	"strconv"

	"todolist/internal/model"
	"todolist/internal/repository"
)

// TodoService 待办事项服务接口
type TodoService interface {
	GetAllTodos() ([]*model.Todo, error)
	GetTodoByID(id string) (*model.Todo, error)
	CreateTodo(req *model.CreateTodoRequest) (*model.Todo, error)
	UpdateTodo(id string, req *model.UpdateTodoRequest) (*model.Todo, error)
	DeleteTodo(id string) error
	ToggleTodo(id string) (*model.Todo, error)
	GetTodosByCompleted(completedStr string) ([]*model.Todo, error)
	GetStats() (*model.TodoStats, error)
}

// todoService 待办事项服务实现
type todoService struct {
	repo repository.TodoRepository
}

// NewTodoService 创建新的待办事项服务
func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{
		repo: repo,
	}
}

func (s *todoService) GetAllTodos() ([]*model.Todo, error) {
	return s.repo.GetAll()
}

func (s *todoService) GetTodoByID(id string) (*model.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *todoService) CreateTodo(req *model.CreateTodoRequest) (*model.Todo, error) {
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}
	
	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}
	
	return todo, nil
}

func (s *todoService) UpdateTodo(id string, req *model.UpdateTodoRequest) (*model.Todo, error) {
	// 先获取现有的todo
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// 更新字段
	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	
	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}
	
	return todo, nil
}

func (s *todoService) DeleteTodo(id string) error {
	return s.repo.Delete(id)
}

func (s *todoService) ToggleTodo(id string) (*model.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	todo.Completed = !todo.Completed
	
	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}
	
	return todo, nil
}

func (s *todoService) GetTodosByCompleted(completedStr string) ([]*model.Todo, error) {
	if completedStr == "" {
		return s.repo.GetAll()
	}
	
	completed, err := strconv.ParseBool(completedStr)
	if err != nil {
		return nil, err
	}
	
	return s.repo.GetByCompleted(completed)
}

func (s *todoService) GetStats() (*model.TodoStats, error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	
	stats := &model.TodoStats{
		Total: len(todos),
	}
	
	for _, todo := range todos {
		if todo.Completed {
			stats.Completed++
		} else {
			stats.Pending++
		}
	}
	
	return stats, nil
}