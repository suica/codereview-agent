package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todolist/internal/model"
	"todolist/internal/repository"
	"todolist/internal/service"
)

// TodoHandler 待办事项处理器
type TodoHandler struct {
	service service.TodoService
}

// NewTodoHandler 创建新的待办事项处理器
func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// GetTodos 获取待办事项列表（支持状态筛选）
func (h *TodoHandler) GetTodos(c *gin.Context) {
	completedStr := c.Query("completed")
	
	todos, err := h.service.GetTodosByCompleted(completedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "completed参数必须是true或false",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":    todos,
		"message": "获取待办事项列表成功",
	})
}

// GetTodoByID 根据ID获取单个待办事项
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	
	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		if err == repository.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "待办事项未找到",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "服务器内部错误",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":    todo,
		"message": "获取待办事项成功",
	})
}

// CreateTodo 创建新的待办事项
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req model.CreateTodoRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}
	
	todo, err := h.service.CreateTodo(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建待办事项失败",
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"data":    todo,
		"message": "待办事项创建成功",
	})
}

// UpdateTodo 更新待办事项
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateTodoRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}
	
	todo, err := h.service.UpdateTodo(id, &req)
	if err != nil {
		if err == repository.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "待办事项未找到",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新待办事项失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":    todo,
		"message": "待办事项更新成功",
	})
}

// DeleteTodo 删除待办事项
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.service.DeleteTodo(id); err != nil {
		if err == repository.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "待办事项未找到",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除待办事项失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "待办事项删除成功",
	})
}

// ToggleTodo 切换待办事项完成状态
func (h *TodoHandler) ToggleTodo(c *gin.Context) {
	id := c.Param("id")
	
	todo, err := h.service.ToggleTodo(id)
	if err != nil {
		if err == repository.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "待办事项未找到",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "切换状态失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":    todo,
		"message": "待办事项状态更新成功",
	})
}

// GetStats 获取统计信息
func (h *TodoHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取统计信息失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":    stats,
		"message": "获取统计信息成功",
	})
}
