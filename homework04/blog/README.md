# 个人博客系统后端

基于 Go + Gin + GORM 开发的个人博客系统后端。

## 功能特性

- **用户管理**
  - 用户注册与登录
  - JWT 认证与授权
  - 用户信息获取

- **文章管理**
  - 文章创建、读取、更新、删除 (CRUD)
  - 文章列表分页查询
  - 文章详情查看

- **评论管理**
  - 评论创建与删除
  - 文章评论列表查询
  - 评论权限控制

- **系统管理**
  - 数据库自动迁移
  - 结构化日志记录
  - 错误处理中间件
  - 健康检查接口

## 技术栈

- **后端框架**: Go 1.21+ + Gin Web 框架
- **数据库**: MySQL + GORM ORM
- **认证**: JWT (JSON Web Token)
- **日志**: Zap 高性能日志库
- **配置管理**: godotenv 环境变量管理
- **密码加密**: bcrypt 密码哈希

## 项目结构
```
blog/
├── config/
│   └── database.go          # 数据库配置和连接，包含环境变量加载和数据库初始化
├── env/
│   └── .env.example         # 环境变量示例文件
├── handlers/
│   ├── auth.go              # 认证相关处理器：注册、登录、获取用户信息
│   ├── comment.go           # 评论相关处理器：创建、获取、删除评论
│   └── post.go              # 文章相关处理器：文章CRUD操作
├── middleware/
│   ├── auth.go              # JWT 认证中间件，验证token有效性
│   └── logger.go            # 日志中间件，记录请求信息
├── models/
│   ├── comment.go           # 评论数据模型，定义评论表结构
│   ├── post.go              # 文章数据模型，定义文章表结构
│   └── user.go              # 用户数据模型，定义用户表结构
├── routes/
│   └── routes.go            # 路由配置，定义所有API端点
├── utils/
│   ├── jwt.go               # JWT 工具函数：生成和验证token
│   └── response.go          # 统一响应格式工具函数
├── .env                     # 环境变量配置文件
├── go.mod                   # Go 模块依赖管理
├── go.sum                   # 依赖校验文件
├── main.go                  # 应用入口文件，初始化数据库、路由和启动服务器
└── README.md                # 项目说明文档
```

### 核心文件说明

- **main.go**: 应用入口，负责初始化配置、数据库、中间件和启动服务器
- **config/database.go**: 数据库连接配置，支持环境变量配置和自动迁移
- **routes/routes.go**: 定义所有API路由，包括认证、文章、评论等模块
- **handlers/**: 业务逻辑处理层，处理具体的HTTP请求
- **models/**: 数据模型定义，对应数据库表结构
- **middleware/**: 中间件层，处理认证、日志等通用功能
- **utils/**: 工具函数，提供JWT和响应格式等通用功能

## API 接口文档

### 认证接口

| 方法 | 路径 | 描述 | 认证要求 |
|------|------|------|----------|
| POST | `/api/auth/register` | 用户注册 | 无需认证 |
| POST | `/api/auth/login` | 用户登录 | 无需认证 |
| GET | `/api/auth/profile` | 获取用户信息 | 需要认证 |

### 文章接口

| 方法 | 路径 | 描述 | 认证要求 |
|------|------|------|----------|
| GET | `/api/posts` | 获取文章列表 | 无需认证 |
| GET | `/api/posts/:id` | 获取文章详情 | 无需认证 |
| POST | `/api/posts` | 创建文章 | 需要认证 |
| PUT | `/api/posts/:id` | 更新文章 | 需要认证 |
| DELETE | `/api/posts/:id` | 删除文章 | 需要认证 |

### 评论接口

| 方法 | 路径 | 描述 | 认证要求 |
|------|------|------|----------|
| GET | `/api/posts/:id/comments` | 获取文章评论 | 无需认证 |
| POST | `/api/posts/:id/comments` | 创建评论 | 需要认证 |
| DELETE | `/api/posts/:id/comments/:commentId` | 删除评论 | 需要认证 |

### 系统接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查 |

## 启动项目

### 环境要求

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本
- Git

### 安装依赖

```bash
# 克隆项目
git clone <repository-url>
cd blog

# 安装 Go 依赖
go mod tidy
```

### 数据库配置

1. 创建 MySQL 数据库
```sql
CREATE DATABASE blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 配置环境变量
复制 `env/.env.example` 为 `.env` 并修改配置：
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=blog

# 服务器配置
PORT=8080
JWT_SECRET=your_jwt_secret_key
```

### 运行项目

```bash
# 开发模式运行
go run main.go

# 或编译后运行
go build -o blog
./blog
```

服务器启动后默认运行在 `http://localhost:8080`

## 接口测试用例和测试结果

### 测试工具
- 使用 Postman 或 curl 进行接口测试
- 测试环境: Windows 11 + Go 1.21 + MySQL 8.0

### 认证接口测试

#### 1. 用户注册
**请求:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com", "password": "password123"}'
```

**预期结果:**
- 状态码: 201 Created
- 响应:
```json
{
  "code": 201,
  "message": "用户注册成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 2. 用户登录
**请求:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

#### 3. 获取用户信息
**请求:**
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "获取用户信息成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 文章接口测试

#### 4. 获取文章列表
**请求:**
```bash
curl -X GET http://localhost:8080/api/posts
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "获取文章列表成功",
  "data": [
    {
      "id": 1,
      "title": "第一篇博客文章",
      "content": "这是第一篇博客文章的内容...",
      "author_id": 1,
      "author_name": "testuser",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "total": 1
}
```

#### 5. 获取文章详情
**请求:**
```bash
curl -X GET http://localhost:8080/api/posts/1
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "获取文章成功",
  "data": {
    "id": 1,
    "title": "第一篇博客文章",
    "content": "这是第一篇博客文章的内容...",
    "author_id": 1,
    "author_name": "testuser",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z",
    "comments_count": 2
  }
}
```

#### 6. 创建文章
**请求:**
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"title": "测试文章", "content": "这是测试文章的内容"}'
```

**预期结果:**
- 状态码: 201 Created
- 响应:
```json
{
  "code": 201,
  "message": "文章创建成功",
  "data": {
    "id": 2,
    "title": "测试文章",
    "content": "这是测试文章的内容",
    "author_id": 1,
    "created_at": "2024-01-01T11:00:00Z",
    "updated_at": "2024-01-01T11:00:00Z"
  }
}
```

#### 7. 更新文章
**请求:**
```bash
curl -X PUT http://localhost:8080/api/posts/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"title": "更新后的测试文章", "content": "这是更新后的内容"}'
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "文章更新成功",
  "data": {
    "id": 2,
    "title": "更新后的测试文章",
    "content": "这是更新后的内容",
    "author_id": 1,
    "created_at": "2024-01-01T11:00:00Z",
    "updated_at": "2024-01-01T11:30:00Z"
  }
}
```

#### 8. 删除文章
**请求:**
```bash
curl -X DELETE http://localhost:8080/api/posts/2 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "文章删除成功",
  "data": null
}
```

### 评论接口测试

#### 9. 获取文章评论
**请求:**
```bash
curl -X GET http://localhost:8080/api/posts/1/comments
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "获取评论列表成功",
  "data": [
    {
      "id": 1,
      "content": "这是一条评论",
      "author_id": 1,
      "author_name": "testuser",
      "post_id": 1,
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 1
}
```

#### 10. 创建评论
**请求:**
```bash
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"content": "这是一条新的评论"}'
```

**预期结果:**
- 状态码: 201 Created
- 响应:
```json
{
  "code": 201,
  "message": "评论创建成功",
  "data": {
    "id": 2,
    "content": "这是一条新的评论",
    "author_id": 1,
    "author_name": "testuser",
    "post_id": 1,
    "created_at": "2024-01-01T13:00:00Z"
  }
}
```

#### 11. 删除评论
**请求:**
```bash
curl -X DELETE http://localhost:8080/api/posts/1/comments/2 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "code": 200,
  "message": "评论删除成功",
  "data": null
}
```

### 系统接口测试

#### 12. 健康检查
**请求:**
```bash
curl -X GET http://localhost:8080/health
```

**预期结果:**
- 状态码: 200 OK
- 响应:
```json
{
  "status": "OK",
  "message": "Blog API is running"
}
```
