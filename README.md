# TodoList 后端服务

一个使用 Go 语言和 Gin 框架实现的 TodoList 后端 API 服务，采用标准的 Go 项目结构和分层架构设计。

## 功能特性

- ✅ 创建待办事项
- ✅ 查看所有待办事项
- ✅ 查看单个待办事项
- ✅ 更新待办事项
- ✅ 删除待办事项
- ✅ 切换完成状态
- ✅ 按状态筛选待办事项
- ✅ 获取统计信息
- ✅ 健康检查
- ✅ CORS 支持

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 运行服务

```bash
go run cmd/todolist/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 健康检查

```
GET /health
```

### 待办事项管理

#### 获取所有待办事项

```
GET /api/v1/todos
```

#### 按状态筛选待办事项

```
GET /api/v1/todos?completed=true   # 获取已完成的
GET /api/v1/todos?completed=false  # 获取未完成的
```

#### 获取单个待办事项

```
GET /api/v1/todos/:id
```

#### 创建待办事项

```
POST /api/v1/todos
Content-Type: application/json

{
  "title": "学习Go语言",
  "description": "完成Go语言基础教程"
}
```

#### 更新待办事项

```
PUT /api/v1/todos/:id
Content-Type: application/json

{
  "title": "更新后的标题",
  "description": "更新后的描述",
  "completed": true
}
```

#### 删除待办事项

```
DELETE /api/v1/todos/:id
```

#### 切换完成状态

```
PATCH /api/v1/todos/:id/toggle
```

#### 获取统计信息

```
GET /api/v1/stats
```

## 数据结构

### Todo 对象

```json
{
  "id": "uuid-string",
  "title": "待办事项标题",
  "description": "待办事项描述",
  "completed": false,
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T10:00:00Z"
}
```

### 响应格式

成功响应：
```json
{
  "data": {},
  "message": "操作成功"
}
```

错误响应：
```json
{
  "error": "错误信息"
}
```

## 示例请求

### 使用 curl 创建待办事项

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "买菜",
    "description": "去超市买今天晚饭的食材"
  }'
```

### 使用 curl 获取所有待办事项

```bash
curl http://localhost:8080/api/v1/todos
```

### 使用 curl 切换完成状态

```bash
curl -X PATCH http://localhost:8080/api/v1/todos/{todo-id}/toggle
```

## 项目结构

```
.
├── cmd/
│   └── todolist/
│       └── main.go          # 应用程序入口
├── internal/
│   ├── handler/
│   │   └── todo.go          # HTTP 处理器层
│   ├── middleware/
│   │   └── cors.go          # 中间件
│   ├── model/
│   │   └── todo.go          # 数据模型
│   ├── repository/
│   │   └── todo.go          # 数据访问层
│   ├── server/
│   │   └── server.go        # 服务器配置
│   └── service/
│       └── todo.go          # 业务逻辑层
├── go.mod
├── go.sum
└── README.md
```

## 架构设计

本项目采用分层架构设计，各层职责清晰：

- **Handler Layer (处理器层)**: 处理 HTTP 请求和响应
- **Service Layer (服务层)**: 处理业务逻辑
- **Repository Layer (仓储层)**: 数据访问抽象
- **Model Layer (模型层)**: 数据结构定义
- **Middleware Layer (中间件层)**: 横切关注点

## 技术栈

- **Go 1.21+**
- **Gin Web Framework** - HTTP 路由和中间件
- **UUID** - 生成唯一标识符

## 注意事项

- 当前使用内存存储数据，服务重启后数据会丢失
- 生产环境建议使用数据库（如 PostgreSQL、MySQL 等）
- 已包含 CORS 支持，可直接被前端应用调用
- 服务启动时会自动创建一些示例数据

## 扩展建议

- 添加数据库支持
- 添加用户认证和授权
- 添加分页功能
- 添加搜索功能
- 添加日志记录
- 添加单元测试
- 添加 Docker 支持