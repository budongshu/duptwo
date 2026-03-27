# DataRegistry 项目架构规范 (SPEC)

> 本文档是 DataRegistry 数据管理平台的架构设计、代码规范和产品逻辑的完整规范文档，可作为后续项目的参考模板。

---

## 1. 项目概述

### 1.1 技术栈

| 层级 | 技术选型 |
|------|----------|
| **后端** | Go 1.21+ / Gin / GORM |
| **前端** | Vue 3 + TypeScript + Vite 5 + Element Plus |
| **数据库** | SQLite (默认) / MySQL / PostgreSQL |
| **认证** | JWT + TOTP (MFA) |
| **部署** | Docker + Docker Compose |

### 1.2 项目结构

```
DataRegistry/
├── backend/                          # Go 后端
│   ├── cmd/server/                   # 入口和静态文件服务
│   │   └── main.go
│   ├── app/
│   │   ├── api/v1/                  # API 处理器层 (Controller)
│   │   │   ├── auth/
│   │   │   ├── user/
│   │   │   ├── role/
│   │   │   ├── user_group/
│   │   │   ├── upload_record/
│   │   │   ├── field_config/
│   │   │   └── audit/
│   │   ├── dto/                     # 数据传输对象
│   │   ├── model/                   # 数据库模型
│   │   ├── repo/                    # 数据访问层
│   │   ├── service/                 # 业务逻辑层
│   │   └── middleware/               # 中间件
│   ├── global/                      # 全局配置和连接
│   ├── router/                      # 路由定义
│   └── docs/                        # Swagger 文档
├── frontend/                        # Vue 前端
│   └── src/
│       ├── api/                     # API 客户端模块
│       ├── components/              # 公共组件
│       ├── composables/             # 组合式函数
│       ├── config/                  # 应用配置
│       ├── enums/                   # TypeScript 枚举
│       ├── lang/                   # 国际化
│       ├── layout/                 # 布局组件
│       ├── routers/                 # 路由配置
│       ├── store/                   # 状态管理
│       ├── styles/                 # 全局样式
│       ├── typings/                 # 类型定义
│       ├── utils/                   # 工具函数
│       └── views/                  # 页面组件
└── deploy/                         # 部署配置
```

---

## 2. 后端架构规范 (Go + Gin)

### 2.1 分层架构

```
HTTP Request
     ↓
[Router] ─── Middleware (JWT / CORS / Permission)
     ↓
[API Layer]    ← 接收请求、参数校验、调用 Service
     ↓
[Service Layer] ← 业务逻辑、事务管理
     ↓
[Repo Layer]   ← 数据库操作
     ↓
[Model/DB]
```

### 2.2 API 层规范

**文件位置**: `backend/app/api/v1/<entity>/`

**命名规范**:
- 文件名: `entity_name.go` (蛇形)
- 结构体名: `<Entity>Api` (PascalCase)
- 构造方法: `New<Entity>Api()` 返回指针

**代码模板**:

```go
// backend/app/api/v1/user/user.go
package api

import (
    "net/http"
    "strconv"

    "datauptwo/app/dto"
    "datauptwo/app/service"
    "datauptwo/middleware"

    "github.com/gin-gonic/gin"
)

type UserApi struct {
    userService  *service.UserService
    auditService *service.AuditService
}

func NewUserApi() *UserApi {
    return &UserApi{
        userService:  service.NewUserService(),
        auditService: service.NewAuditService(),
    }
}

// Create 创建用户
// @Summary 创建用户
// @Tags User
// @Accept json
// @Param request body dto.UserCreateReq true "用户信息"
// @Success 200 {object} dto.Response
// @Router /api/users [post]
func (api *UserApi) Create(c *gin.Context) {
    var req dto.UserCreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
        return
    }

    user, err := api.userService.Create(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
        return
    }

    // 记录审计日志
    api.auditService.Create(c, "create", "user", user.ID, user.Username)

    c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: user})
}

// GetByID 获取用户详情
func (api *UserApi) GetByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    user, err := api.userService.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "用户不存在"})
        return
    }

    c.JSON(http.StatusOK, dto.Response{Code: 200, Data: user})
}

// Delete 删除用户 (软删除)
func (api *UserApi) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    if err := api.userService.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
        return
    }

    api.auditService.Create(c, "delete", "user", uint(id), "")
    c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除
func (api *UserApi) BatchDelete(c *gin.Context) {
    var req dto.BatchDeleteReq
    if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
        c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的记录"})
        return
    }

    if err := api.userService.BatchDelete(req.IDs); err != nil {
        c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
        return
    }

    api.auditService.Create(c, "batch_delete", "user", 0, fmt.Sprintf("删除了 %d 个用户", len(req.IDs)))
    c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功删除 %d 条记录", len(req.IDs))})
}

// ============ 私有辅助方法 ============
func (api *UserApi) getUserID(c *gin.Context) uint {
    id, _ := c.Get("userId")
    return id.(uint)
}
```

### 2.3 Service 层规范

**文件位置**: `backend/app/service/`

**命名规范**:
- 文件名: `entity_service.go` 或 `entity_name.go`
- 结构体名: `<Entity>Service`
- 构造方法: `New<Entity>Service()`

**代码模板**:

```go
// backend/app/service/user_service.go
package service

type UserService struct {
    userRepo  *repo.UserRepo
    roleRepo  *repo.RoleRepo
    groupRepo *repo.UserGroupRepo
}

func NewUserService() *UserService {
    return &UserService{
        userRepo:  repo.NewUserRepo(),
        roleRepo:  repo.NewRoleRepo(),
        groupRepo: repo.NewUserGroupRepo(),
    }
}

func (s *UserService) Create(req dto.UserCreateReq) (*model.User, error) {
    // 1. 校验逻辑
    if req.Username == "" {
        return nil, errors.New("用户名不能为空")
    }

    // 2. 检查唯一性
    exists, _ := s.userRepo.ExistsByUsername(req.Username)
    if exists {
        return nil, errors.New("用户名已存在")
    }

    // 3. 密码加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // 4. 构建模型
    user := &model.User{
        Username: req.Username,
        Password: string(hashedPassword),
        Email:    req.Email,
        Nickname: req.Nickname,
        Status:   "active",
    }

    // 5. 保存
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

// BatchDelete 批量删除
func (s *UserService) BatchDelete(ids []uint) error {
    if len(ids) == 0 {
        return nil
    }
    return s.userRepo.BatchDelete(ids)
}

// toUserResp 转换为响应 DTO
func (s *UserService) toUserResp(user *model.User) *dto.UserResp {
    if user == nil {
        return nil
    }
    return &dto.UserResp{
        ID:         user.ID,
        Username:   user.Username,
        Nickname:   user.Nickname,
        Email:      user.Email,
        Status:     user.Status,
        RoleID:     user.RoleID,
        RoleName:   "", // 可从 roleRepo 获取
        CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
    }
}
```

### 2.4 Repo 层规范

**文件位置**: `backend/app/repo/`

**命名规范**:
- 文件名: `entity_repo.go`
- 结构体名: `<Entity>Repo`
- 构造方法: `New<Entity>Repo()`

**代码模板**:

```go
// backend/app/repo/user_repo.go
package repo

import (
    "datauptwo/app/dto"
    "datauptwo/app/model"
    "datauptwo/global"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
    return &UserRepo{}
}

// Create 创建记录
func (r *UserRepo) Create(user *model.User) error {
    return global.DB.Create(user).Error
}

// Update 更新记录
func (r *UserRepo) Update(user *model.User) error {
    return global.DB.Save(user).Error
}

// Delete 软删除 (设置 is_deleted = true)
func (r *UserRepo) Delete(id uint) error {
    return global.DB.Model(&model.User{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// BatchDelete 批量软删除
func (r *UserRepo) BatchDelete(ids []uint) error {
    if len(ids) == 0 {
        return nil
    }
    return global.DB.Model(&model.User{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

// GetByID 获取单条记录 (排除已删除)
func (r *UserRepo) GetByID(id uint) (*model.User, error) {
    var user model.User
    err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
    var count int64
    err := global.DB.Model(&model.User{}).
        Where("username = ? AND is_deleted = ?", username, false).
        Count(&count).Error
    return count > 0, err
}

// List 分页列表查询
func (r *UserRepo) List(req dto.UserListReq) ([]model.User, int64, error) {
    var users []model.User
    var total int64

    db := global.DB.Model(&model.User{}).Where("is_deleted = ?", false)

    // 动态条件
    if req.Keyword != "" {
        db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
            "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
    }
    if req.Status != "" {
        db = db.Where("status = ?", req.Status)
    }
    if req.RoleID > 0 {
        db = db.Where("role_id = ?", req.RoleID)
    }

    // 计数
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // 分页查询
    offset := (req.Page - 1) * req.PageSize
    if err := db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}
```

### 2.5 Model 层规范

**文件位置**: `backend/app/model/`

```go
// base.go - 基础模型
package model

import (
    "time"
    "gorm.io/gorm"
)

// 基础模型 (所有模型嵌入此结构)
type BaseModel struct {
    ID        uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// 用户模型
type User struct {
    BaseModel
    Username    string  `gorm:"size:64;not null;uniqueIndex"`
    Password    string  `gorm:"size:256;not null"`  // bcrypt 加密存储
    Nickname    string  `gorm:"size:64"`
    Email       string  `gorm:"size:128;index"`
    Phone       string  `gorm:"size:32"`
    Avatar      string  `gorm:"size:512"`
    Status      string  `gorm:"size:16;default:active"` // active / inactive
    RoleID      uint    `gorm:"index"`
    GroupID     uint    `gorm:"index"`
    LastLoginAt *time.Time
    LastLoginIP string `gorm:"size:64"`
    IsDeleted   bool    `gorm:"default:false;index"` // 软删除标记
    MFAEnabled  bool    `gorm:"default:false"`
    MFASecret   string  `gorm:"size:128"`
}

// 角色模型
type Role struct {
    BaseModel
    Name        string `gorm:"size:64;not null"`
    Code        string `gorm:"size:64;uniqueIndex;not null"`
    Description string `gorm:"size:256"`
    Permissions string `gorm:"type:text"` // JSON 数组格式
    IsDeleted   bool   `gorm:"default:false"`
    Sort        int    `gorm:"default:0"`
}

// 业务模型示例
type UploadRecord struct {
    BaseModel
    SerialNo    string `gorm:"size:64;index"`
    DataType    string `gorm:"size:64;index"`
    ProjectName string `gorm:"size:128;index"`
    FilePath    string `gorm:"size:512"`
    FileSize    int64  `gorm:"default:0"`
    Uploader    string `gorm:"size:64;index"`
    Status      string `gorm:"size:16;default:pending"` // pending/processing/completed/failed
    Remark      string `gorm:"size:512"`
    Data        string `gorm:"type:text"` // JSON 动态字段
    IsDeleted   bool   `gorm:"default:false"`
}
```

### 2.6 DTO 层规范

**文件位置**: `backend/app/dto/`

```go
// common.go - 通用 DTO
package dto

// 统一响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 分页结果
type PageResult struct {
    Total int64       `json:"total"`
    Items interface{} `json:"items"`
}

// 通用分页请求
type PageReq struct {
    Page     int    `json:"page" form:"page" binding:"required,min=1"`
    PageSize int    `json:"pageSize" form:"pageSize" binding:"required,min=1,max=100"`
    Keyword  string `json:"keyword" form:"keyword"`
}

// 批量删除请求
type BatchDeleteReq struct {
    IDs []uint `json:"ids" binding:"required,min=1"`
}

// user.go - 用户相关 DTO
package dto

type UserCreateReq struct {
    Username string `json:"username" binding:"required,min=3,max=32"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"omitempty,email"`
    Nickname string `json:"nickname"`
    Phone    string `json:"phone"`
    RoleID   uint   `json:"roleId"`
    GroupID  uint   `json:"groupId"`
}

type UserUpdateReq struct {
    ID       uint   `json:"id" binding:"required"`
    Nickname string `json:"nickname"`
    Email    string `json:"email" binding:"omitempty,email"`
    Phone    string `json:"phone"`
    RoleID   uint   `json:"roleId"`
    GroupID  uint   `json:"groupId"`
    Status   string `json:"status"` // active / inactive
}

type UserListReq struct {
    PageReq
    Status  string `form:"status"`
    RoleID  uint   `form:"roleId"`
    GroupID uint   `form:"groupId"`
}

type UserResp struct {
    ID         uint   `json:"id"`
    Username   string `json:"username"`
    Nickname   string `json:"nickname"`
    Email      string `json:"email"`
    Phone      string `json:"phone"`
    Status     string `json:"status"`
    RoleID     uint   `json:"roleId"`
    RoleName   string `json:"roleName"`
    GroupID    uint   `json:"groupId"`
    GroupName  string `json:"groupName"`
    MFAEnabled bool   `json:"mfaEnabled"`
    CreatedAt  string `json:"createdAt"`
}
```

### 2.7 路由注册规范

**文件位置**: `backend/router/router.go`

```go
package router

func NewRouter() *gin.Engine {
    r := gin.Default()

    // CORS 中间件
    r.Use(middleware.CORS())

    // 公开接口 (无需认证)
    public := r.Group("/public")
    {
        public.POST("/upload-records", RouterGroupApp.PublicUploadRecordApi.Create)
        public.GET("/upload-records/:serialNo", RouterGroupApp.PublicUploadRecordApi.GetBySerialNo)
    }

    // 认证接口 (无需 JWT)
    r.GET("/api/auth/captcha", RouterGroupApp.AuthApi.GetCaptcha)
    r.POST("/api/auth/login", RouterGroupApp.AuthApi.Login)
    r.POST("/api/auth/register", RouterGroupApp.AuthApi.Register)

    // 受保护接口 (需要 JWT + 权限)
    authGroup := r.Group("/api")
    authGroup.Use(middleware.JWTAuth())
    {
        // === 上传记录 ===
        authGroup.GET("/upload-records", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.List)
        authGroup.POST("/upload-records", middleware.RequirePermission("upload:create"), RouterGroupApp.UploadRecordApi.Create)
        authGroup.POST("/upload-records/batch-delete", middleware.RequirePermission("upload:delete"), RouterGroupApp.UploadRecordApi.BatchDelete)
        authGroup.DELETE("/upload-records/:id", middleware.RequirePermission("upload:delete"), RouterGroupApp.UploadRecordApi.Delete)

        // === 用户管理 ===
        authGroup.GET("/users", middleware.RequirePermission("user:read"), RouterGroupApp.UserApi.List)
        authGroup.POST("/users", middleware.RequirePermission("user:create"), RouterGroupApp.UserApi.Create)
        authGroup.PUT("/users", middleware.RequirePermission("user:update"), RouterGroupApp.UserApi.Update)
        authGroup.DELETE("/users/:id", middleware.RequirePermission("user:delete"), RouterGroupApp.UserApi.Delete)
        authGroup.POST("/users/batch-delete", middleware.RequirePermission("user:delete"), RouterGroupApp.UserApi.BatchDelete)
    }

    return r
}
```

### 2.8 中间件规范

**JWT 认证中间件**:

```go
// middleware/jwt.go
package middleware

type Claims struct {
    UserID      uint     `json:"userId"`
    Username    string   `json:"username"`
    IsAdmin     bool     `json:"isAdmin,omitempty"`
    Permissions []string `json:"permissions,omitempty"`
    jwt.RegisteredClaims
}

// JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := GetToken(c) // 从 Header 或 Cookie 获取
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "message": "未登录"})
            return
        }

        claims, err := ParseToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "message": "登录已过期"})
            return
        }

        // 将用户信息存入 Context
        c.Set("userId", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("isAdmin", claims.IsAdmin)
        c.Set("permissions", claims.Permissions)

        c.Next()
    }
}

// 权限检查中间件
func RequirePermission(perm string) gin.HandlerFunc {
    return func(c *gin.Context) {
        permissions, _ := c.Get("permissions")
        perms := permissions.([]string)

        if !contains(perms, perm) {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无此权限: " + perm})
            c.Abort()
            return
        }
        c.Next()
    }
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

### 2.9 权限码规范

```json
[
  "upload:create", "upload:read", "upload:update", "upload:delete", "upload:export",
  "field-config:create", "field-config:read", "field-config:update", "field-config:delete",
  "user:create", "user:read", "user:update", "user:delete",
  "role:create", "role:read", "role:update", "role:delete",
  "audit:operation:read", "audit:login:read",
  "config:read", "config:update",
  "admin:all"
]
```

**默认角色**:
- `admin`: 拥有 `admin:all` 和所有权限
- `auditor`: 审计相关只读权限
- `readonly`: 仅查看上传记录

---

## 3. 前端架构规范 (Vue 3 + TypeScript)

### 3.1 项目结构

```
frontend/src/
├── api/                    # API 客户端模块 (按模块划分)
│   ├── index.ts           # Axios 实例和通用响应类型
│   ├── auth.ts            # 认证相关 API
│   ├── user.ts            # 用户管理 API
│   ├── role.ts            # 角色管理 API
│   └── upload-record.ts   # 上传记录 API
├── components/            # 公共组件
│   └── ...               # 业务无关的通用组件
├── composables/           # 组合式函数 (复用逻辑)
├── config/               # 应用配置
├── enums/                # TypeScript 枚举
├── lang/                 # 国际化文件
├── layout/               # 布局组件
│   └── index.vue         # 侧边栏 + 主内容布局
├── routers/             # 路由配置
│   ├── index.ts          # 路由实例
│   ├── static.ts         # 静态路由 (公开)
│   └── async.ts          # 动态路由 (需权限)
├── store/                # Pinia 状态管理
├── styles/               # 全局 SCSS 样式
├── typings/              # 全局类型定义
├── utils/                # 工具函数
└── views/                # 页面组件 (按模块组织)
    ├── auth/
    │   ├── login.vue
    │   └── register.vue
    ├── user/
    │   └── list.vue      # 用户列表页
    ├── role/
    │   └── list.vue      # 角色列表页
    └── upload-record/
        ├── dashboard.vue # 数据看板
        └── list.vue      # 上传记录列表
```

### 3.2 API 客户端规范

**Axios 实例配置** (`api/index.ts`):

```typescript
// api/index.ts
import axios, { AxiosInstance, AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

// ============ 通用响应类型 ============
export interface ResultData<T = any> {
  code: number
  message: string
  data: T
}

export interface PageResult<T = any> {
  total: number
  items: T[]
}

// ============ 请求封装 ============
class RequestHttp {
  service: AxiosInstance

  constructor(cfg: AxiosRequestConfig) {
    this.service = axios.create({
      baseURL: import.meta.env.VITE_API_URL || '/api',
      timeout: 30000,
      ...cfg,
    })

    // 请求拦截器: 添加 Token
    this.service.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    // 响应拦截器: 处理错误
    this.service.interceptors.response.use(
      (response: AxiosResponse<ResultData>) => {
        // 文件流直接返回
        if (response.config.responseType === 'blob') {
          return response.data
        }

        const { data } = response
        if (data.code !== 200) {
          ElMessage.error(data.message || '请求失败')
          return Promise.reject(data)
        }
        return data
      },
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token')
          window.location.href = '/login'
        }
        ElMessage.error(error.response?.data?.message || '网络错误')
        return Promise.reject(error.response?.data || error)
      }
    )
  }

  get<T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<ResultData<T>> {
    return this.service.get(url, { params, ...config })
  }

  post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ResultData<T>> {
    return this.service.post(url, data, config)
  }

  put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ResultData<T>> {
    return this.service.put(url, data, config)
  }

  delete<T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<ResultData<T>> {
    return this.service.delete(url, { params, ...config })
  }
}

export default new RequestHttp({})
```

**API 模块规范** (`api/user.ts`):

```typescript
// api/user.ts
import request from './index'
import type { PageResult } from './index'

// ============ 类型定义 ============
export interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone?: string
  status: 'active' | 'inactive'
  roleId: number
  roleName: string
  groupId: number
  groupName: string
  mfaEnabled: boolean
  createdAt: string
  lastLoginAt?: string
}

export interface CreateUserReq {
  username: string
  password: string
  nickname?: string
  email: string
  phone?: string
  roleId?: number
  groupId?: number
}

export interface UpdateUserReq {
  id: number
  nickname?: string
  email?: string
  phone?: string
  roleId?: number
  groupId?: number
  status?: 'active' | 'inactive'
}

export interface UserListReq {
  page?: number
  pageSize?: number
  keyword?: string
  status?: string
  roleId?: number
  groupId?: number
}

// ============ API 命名空间 ============
export namespace UserApi {
  // 获取用户列表
  export const list = (params?: UserListReq) => {
    return request.get<PageResult<User[]>>('/users', params)
  }

  // 获取用户详情
  export const getById = (id: number) => {
    return request.get<User>(`/users/${id}`)
  }

  // 创建用户
  export const create = (data: CreateUserReq) => {
    return request.post<User>('/users', data)
  }

  // 更新用户
  export const update = (data: UpdateUserReq) => {
    return request.put<User>('/users', data)
  }

  // 删除用户
  export const del = (id: number) => {
    return request.delete<null>(`/users/${id}`)
  }

  // 批量删除
  export const batchDelete = (ids: number[]) => {
    return request.post<null>('/users/batch-delete', { ids })
  }
}
```

### 3.3 列表页规范 (CRUD)

**页面结构**:

```vue
<!-- views/user/list.vue -->
<template>
  <div class="page">
    <!-- 页面标题栏 -->
    <header class="page-header">
      <div>
        <h1 class="page-title">用户管理</h1>
        <p class="page-subtitle">管理系统用户账号</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <IconPlus />
          新增用户
        </el-button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-card">
      <el-input v-model="keyword" placeholder="搜索用户名/邮箱" clearable @keyup.enter="handleSearch" />
      <el-select v-model="status" placeholder="状态" clearable>
        <el-option label="正常" value="active" />
        <el-option label="禁用" value="inactive" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 表格卡片 -->
    <div class="table-card">
      <!-- 表格工具栏 -->
      <div class="table-toolbar">
        <div class="toolbar-left">
          <span class="record-count">共 <strong>{{ pagination.total }}</strong> 条</span>
          <span v-if="selectedRows.length > 0" class="selection-count">
            已选 <strong>{{ selectedRows.length }}</strong> 项
          </span>
        </div>
        <div class="toolbar-right">
          <!-- 批量删除按钮 -->
          <el-button
            v-if="selectedRows.length > 0"
            type="danger"
            @click="handleBatchDelete"
          >
            <IconDelete />
            批量删除
          </el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table
        ref="tableRef"
        v-model:selection="selectedRows"
        :data="tableData"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="45" fixed="left" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="100" />
        <el-table-column prop="email" label="邮箱" min-width="160" />
        <el-table-column prop="roleName" label="角色" min-width="100" align="center" />
        <el-table-column prop="status" label="状态" min-width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          background
        />
      </div>
    </div>

    <!-- 编辑抽屉 -->
    <el-drawer v-model="drawerVisible" :title="isEdit ? '编辑用户' : '新增用户'" size="480px" direction="rtl">
      <div class="drawer-content">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
          <el-form-item label="用户名" prop="username" v-if="!isEdit">
            <el-input v-model="form.username" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="密码" prop="password" v-if="!isEdit">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="form.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="角色">
            <el-select v-model="form.roleId" placeholder="请选择角色" style="width: 100%">
              <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="drawer-footer">
          <el-button @click="drawerVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="confirmSubmit">保存</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UserApi, type User, type CreateUserReq, type UpdateUserReq } from '@/api/user'

// ============ 状态 ============
const loading = ref(false)
const submitting = ref(false)
const tableData = ref<User[]>([])
const tableRef = ref()
const selectedRows = ref<User[]>([])
const drawerVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const roles = ref<Role[]>([])

const keyword = ref('')
const status = ref('')
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive<CreateUserReq & { id?: number }>({
  username: '',
  password: '',
  nickname: '',
  email: '',
  roleId: undefined,
})

const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '用户名长度为 3-32 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
}

// ============ 方法 ============
const loadData = async () => {
  loading.value = true
  try {
    const res = await UserApi.list({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: keyword.value || undefined,
      status: status.value || undefined,
    })
    if (res.code === 200) {
      tableData.value = res.data.items || []
      pagination.total = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  keyword.value = ''
  status.value = ''
  handleSearch()
}

const handleSelectionChange = (rows: User[]) => {
  selectedRows.value = rows
}

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { id: undefined, username: '', password: '', nickname: '', email: '', roleId: undefined })
  drawerVisible.value = true
}

const handleEdit = (row: User) => {
  isEdit.value = true
  Object.assign(form, { id: row.id, username: row.username, password: '', nickname: row.nickname, email: row.email, roleId: row.roleId })
  drawerVisible.value = true
}

const confirmSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      const data: UpdateUserReq = { id: form.id!, nickname: form.nickname, email: form.email, roleId: form.roleId }
      const res = await UserApi.update(data)
      if (res.code === 200) {
        ElMessage.success('更新成功')
        drawerVisible.value = false
        loadData()
      }
    } else {
      const res = await UserApi.create(form as CreateUserReq)
      if (res.code === 200) {
        ElMessage.success('创建成功')
        drawerVisible.value = false
        loadData()
      }
    }
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: User) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户 "${row.username}" 吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    const res = await UserApi.del(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      loadData()
    }
  } catch {}
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个用户吗？`,
      '批量删除确认',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    const ids = selectedRows.value.map((row) => row.id)
    await UserApi.batchDelete(ids)
    ElMessage.success(`成功删除 ${selectedRows.value.length} 个用户`)
    selectedRows.value = []
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

onMounted(() => loadData())
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f5f6fa;
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 20px 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.page-title {
  font-size: 20px;
  font-weight: 800;
  color: #191c23;
}

.page-subtitle {
  font-size: 13px;
  color: #727785;
}

.filter-card,
.table-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 14px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.filter-card {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;

  .el-input,
  .el-select {
    width: 150px;
  }
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #eaecf0;
  background: #f8f9fc;
}

.toolbar-left {
  .record-count {
    font-size: 13px;
    color: #727785;
    strong {
      color: #191c23;
      font-weight: 700;
    }
  }
  .selection-count {
    margin-left: 16px;
    font-size: 13px;
    color: #c0392b;
    strong {
      font-weight: 700;
    }
  }
}

.toolbar-right {
  display: flex;
  gap: 10px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid #f0f2f5;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 24px;
  border-top: 1px solid #f0f2f5;
}
</style>
```

### 3.4 登录页面规范

```vue
<!-- views/auth/login.vue -->
<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="title">DataRegistry</h1>
      <p class="subtitle">数据管理平台</p>

      <el-form ref="formRef" :model="form" :rules="rules" size="large">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>

        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>

        <!-- 验证码 -->
        <el-form-item prop="captcha">
          <el-input v-model="form.captcha" placeholder="验证码" style="width: 60%" />
          <el-image :src="captchaUrl" class="captcha-img" @click="refreshCaptcha">
            <template #error>
              <div class="captcha-error" @click="refreshCaptcha">点击刷新</div>
            </template>
          </el-image>
        </el-form-item>

        <!-- MFA (可选) -->
        <el-form-item v-if="mfaRequired" prop="mfaCode">
          <el-input v-model="form.mfaCode" placeholder="请输入 MFA 验证码" maxlength="6" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%" @click="handleLogin">
            {{ mfaRequired ? '验证' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { AuthApi } from '@/api/auth'

const router = useRouter()
const loading = ref(false)
const mfaRequired = ref(false)
const captchaUrl = ref('')
const formRef = ref()

const form = reactive({
  username: '',
  password: '',
  captcha: '',
  captchaId: '',
  mfaCode: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

const refreshCaptcha = async () => {
  const res = await AuthApi.getCaptcha()
  if (res.code === 200) {
    form.captchaId = res.data.id
    captchaUrl.value = res.data.image
  }
}

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const loginData = mfaRequired.value
      ? { username: form.username, password: form.password, mfaCode: form.mfaCode }
      : { username: form.username, password: form.password, captcha: form.captcha, captchaId: form.captchaId }

    const res = mfaRequired.value
      ? await AuthApi.mfaVerify({ userId: res.data.userId, code: form.mfaCode })
      : await AuthApi.login(loginData as any)

    if (res.code === 200) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      ElMessage.success('登录成功')
      router.push('/')
    }
  } catch (error: any) {
    if (error.mfaRequired) {
      mfaRequired.value = true
    }
  } finally {
    loading.value = false
  }
}

// 初始化验证码
refreshCaptcha()
</script>
```

### 3.5 路由规范

```typescript
// routers/static.ts - 公开路由
export const staticRoutes: RouteRecordRaw[] = [
  { path: '/', redirect: '/upload-record' },
  { path: '/login', name: 'Login', component: () => import('@/views/auth/login.vue') },
  { path: '/register', name: 'Register', component: () => import('@/views/auth/register.vue') },
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('@/views/error/404.vue') },
]

// routers/async.ts - 需要权限的路由
export const asyncRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    children: [
      { path: 'upload-record', name: 'UploadRecord', component: () => import('@/views/upload-record/dashboard.vue') },
      { path: 'upload-record/list', name: 'UploadRecordList', component: () => import('@/views/upload-record/list.vue') },
      { path: 'users', name: 'UserList', component: () => import('@/views/user/list.vue') },
      { path: 'roles', name: 'RoleList', component: () => import('@/views/role/list.vue') },
      { path: 'user-groups', name: 'UserGroupList', component: () => import('@/views/user-group/list.vue') },
      { path: 'field-config', name: 'FieldConfig', component: () => import('@/views/field-config/list.vue') },
      { path: 'audit/operation-log', name: 'OperationLog', component: () => import('@/views/audit/operation-log.vue') },
      { path: 'audit/login-log', name: 'LoginLog', component: () => import('@/views/audit/login-log.vue') },
    ],
  },
]
```

### 3.6 布局组件规范

```vue
<!-- layout/index.vue -->
<template>
  <div class="app-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed }">
      <div class="logo">
        <img src="@/assets/logo.png" alt="Logo" />
        <span v-if="!collapsed">DataRegistry</span>
      </div>

      <el-menu :default-active="route.path" :collapse="collapsed" router>
        <el-sub-menu index="1">
          <template #title><IconData />数据管理</template>
          <el-menu-item index="/upload-record">上传记录</el-menu-item>
          <el-menu-item index="/field-config">字段配置</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="2">
          <template #title><IconSetting />系统管理</template>
          <el-menu-item index="/users">用户管理</el-menu-item>
          <el-menu-item index="/roles">角色管理</el-menu-item>
          <el-menu-item index="/user-groups">用户组</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="3">
          <template #title><IconDocument />审计日志</template>
          <el-menu-item index="/audit/operation-log">操作日志</el-menu-item>
          <el-menu-item index="/audit/login-log">登录日志</el-menu-item>
        </el-sub-menu>
      </el-menu>

      <div class="sidebar-footer">
        <div class="user-info">
          <el-avatar :size="32">{{ user?.username?.[0] }}</el-avatar>
          <span v-if="!collapsed">{{ user?.username }}</span>
        </div>
        <el-button link @click="handleLogout">退出</el-button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const collapsed = ref(false)
const user = computed(() => {
  const userStr = localStorage.getItem('user')
  return userStr ? JSON.parse(userStr) : null
})

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style scoped lang="scss">
.app-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 220px;
  background: #fff;
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;

  &.collapsed {
    width: 64px;
  }
}

.main-content {
  flex: 1;
  background: #f5f6fa;
  overflow-y: auto;
}
</style>
```

---

## 4. 数据库设计规范

### 4.1 通用字段

```sql
-- 所有表都包含以下基础字段
id          INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
is_deleted  TINYINT(1) DEFAULT 0  -- 软删除标记
```

### 4.2 表结构示例

```sql
-- 用户表
CREATE TABLE users (
    id          INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username    VARCHAR(64) NOT NULL UNIQUE,
    password    VARCHAR(256) NOT NULL,
    nickname    VARCHAR(64),
    email       VARCHAR(128),
    phone       VARCHAR(32),
    avatar      VARCHAR(512),
    status      VARCHAR(16) DEFAULT 'active',
    role_id     INT UNSIGNED,
    group_id    INT UNSIGNED,
    last_login_at DATETIME,
    last_login_ip VARCHAR(64),
    is_deleted  TINYINT(1) DEFAULT 0,
    mfa_enabled TINYINT(1) DEFAULT 0,
    mfa_secret  VARCHAR(128),
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_role (role_id),
    INDEX idx_deleted (is_deleted)
);

-- 角色表
CREATE TABLE roles (
    id          INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(64) NOT NULL,
    code        VARCHAR(64) NOT NULL UNIQUE,
    description VARCHAR(256),
    permissions TEXT,  -- JSON 数组
    sort        INT DEFAULT 0,
    is_deleted  TINYINT(1) DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 上传记录表
CREATE TABLE upload_records (
    id          INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    serial_no   VARCHAR(64) NOT NULL UNIQUE,
    data_type   VARCHAR(64),
    project_name VARCHAR(128),
    file_path   VARCHAR(512),
    file_size   BIGINT DEFAULT 0,
    uploader    VARCHAR(64),
    status      VARCHAR(16) DEFAULT 'pending',
    remark      VARCHAR(512),
    data        TEXT,  -- JSON 动态字段
    is_deleted  TINYINT(1) DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_serial (serial_no),
    INDEX idx_type (data_type),
    INDEX idx_uploader (uploader),
    INDEX idx_created (created_at)
);
```

---

## 5. API 设计规范

### 5.1 RESTful 规范

| 操作 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 列表 | GET | `/api/users` | 分页查询 |
| 详情 | GET | `/api/users/:id` | 获取单条 |
| 创建 | POST | `/api/users` | 新增 |
| 更新 | PUT | `/api/users` | 更新 |
| 删除 | DELETE | `/api/users/:id` | 删除 |
| 批量 | POST | `/api/users/batch-delete` | 批量删除 |

### 5.2 响应格式

```json
// 成功响应
{
  "code": 200,
  "message": "操作成功",
  "data": { ... }
}

// 分页响应
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "items": [...]
  }
}

// 错误响应
{
  "code": 400,
  "message": "参数错误：用户名不能为空",
  "data": null
}
```

### 5.3 HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录或 Token 过期 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 6. 安全规范

### 6.1 认证

- 密码使用 bcrypt 加密存储
- JWT Token 有效期 24 小时
- 支持 TOTP MFA 二次验证
- 敏感操作需要重新验证

### 6.2 权限控制

- 基于角色的权限控制 (RBAC)
- 接口级别权限检查
- 前端路由权限控制 (可选)

### 6.3 输入校验

- 后端所有输入必须校验
- 使用 DTO binding 进行类型转换
- SQL 注入防护 (使用 GORM 参数化查询)

---

## 7. 代码命名规范

### 7.1 后端 (Go)

| 类型 | 规范 | 示例 |
|------|------|------|
| 变量 | camelCase | `userName`, `pageSize` |
| 常量 | PascalCase | `StatusActive`, `RoleAdmin` |
| 函数 | PascalCase | `GetUserByID`, `CreateUser` |
| 结构体 | PascalCase | `UserService`, `UserRepo` |
| 接口 | PascalCase | `UserApi`, `AuthService` |
| 包名 | 全小写蛇形 | `upload_record`, `user_group` |
| 数据库表 | 全小写蛇形复数 | `users`, `upload_records` |
| JSON 字段 | camelCase | `userName`, `createdAt` |

### 7.2 前端 (TypeScript/Vue)

| 类型 | 规范 | 示例 |
|------|------|------|
| 变量/函数 | camelCase | `userName`, `handleClick` |
| 类/接口/类型 | PascalCase | `User`, `UserListReq` |
| 组件名 | PascalCase | `UserList`, `UploadRecord` |
| 常量 | PascalCase | `StatusEnum`, `ApiUrls` |
| 文件名 | kebab-case | `user-list.vue`, `upload-record.ts` |
| CSS 类名 | kebab-case | `page-header`, `table-card` |

---

## 8. 错误处理规范

### 8.1 后端

```go
// 标准错误响应
c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "错误描述"})

// 带数据的成功响应
c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "成功", Data: user})

// 错误日志记录
logger.Error("删除失败", zap.Error(err), zap.Uint("id", id))
```

### 8.2 前端

```typescript
// API 错误处理
try {
  const res = await UserApi.create(data)
  if (res.code === 200) {
    ElMessage.success('创建成功')
    loadData()
  }
} catch (error: any) {
  // 业务错误已由拦截器处理
  // 此处可处理特殊逻辑
}

// 确认对话框
await ElMessageBox.confirm('确定要删除吗？', '提示', { type: 'warning' })

// 表单验证
const valid = await formRef.value?.validate().catch(() => false)
if (!valid) return
```

---

## 9. 项目启动命令

### 9.1 后端

```bash
cd backend

# 开发模式
go run cmd/server/main.go

# 构建
go build -o server ./cmd/server/

# SQLite 开发配置 (默认)
# 监听端口: 18421
```

### 9.2 前端

```bash
cd frontend

# 安装依赖
npm install

# 开发服务器 (端口 4004，代理到后端 18421)
npm run dev

# 生产构建
npm run build
# 输出到 ../backend/cmd/server/web/
```

### 9.3 Docker 部署

```bash
cd deploy

# SQLite 版本
docker-compose -f docker-compose-sqlite.yml up -d

# MySQL 版本
docker-compose up -d
```

---

## 10. 项目检查清单

新增功能时，确保完成以下所有部分：

- [ ] **Model**: 在 `backend/app/model/` 添加数据模型
- [ ] **DTO**: 在 `backend/app/dto/` 添加请求/响应 DTO
- [ ] **Repo**: 在 `backend/app/repo/` 实现数据访问层
- [ ] **Service**: 在 `backend/app/service/` 实现业务逻辑
- [ ] **API**: 在 `backend/app/api/v1/` 实现 HTTP 处理器
- [ ] **Router**: 在 `backend/router/router.go` 注册路由和权限
- [ ] **Frontend API**: 在 `frontend/src/api/` 添加 API 模块
- [ ] **Frontend Page**: 在 `frontend/src/views/` 添加页面组件
- [ ] **Frontend Router**: 在 `frontend/src/routers/` 注册路由
- [ ] **权限码**: 更新默认角色的权限配置
- [ ] **API 文档**: 更新 Swagger 注解
- [ ] **测试**: 编写单元测试或集成测试

---

*本文档由 Claude Code 自动生成，可根据实际项目需求进行调整。*
