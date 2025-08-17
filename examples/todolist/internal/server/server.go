package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"todolist/internal/handler"
	"todolist/internal/middleware"
	"todolist/internal/service"
	"todolist/internal/repository"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	// 初始化仓储层
	todoRepo := repository.NewMemoryTodoRepository()
	
	// 初始化服务层
	todoService := service.NewTodoService(todoRepo)
	
	// 初始化处理器层
	todoHandler := handler.NewTodoHandler(todoService)
	
	// 创建路由器
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	
	// 设置路由
	setupRoutes(router, todoHandler)
	
	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	return s.router.Run(addr)
}

func setupRoutes(r *gin.Engine, todoHandler *handler.TodoHandler) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "TodoList服务运行正常",
			"time":    time.Now(),
		})
	})
	
	// API路由组
	api := r.Group("/api/v1")
	{
		// 待办事项相关路由
		api.GET("/todos", todoHandler.GetTodos)
		api.GET("/todos/:id", todoHandler.GetTodoByID)
		api.POST("/todos", todoHandler.CreateTodo)
		api.PUT("/todos/:id", todoHandler.UpdateTodo)
		api.DELETE("/todos/:id", todoHandler.DeleteTodo)
		api.PATCH("/todos/:id/toggle", todoHandler.ToggleTodo)
		
		// 统计信息
		api.GET("/stats", todoHandler.GetStats)
	}
}