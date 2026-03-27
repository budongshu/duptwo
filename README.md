# DataRegistry - 数据登记管理平台

## 项目简介

一套完整的数据登记管理解决方案，通过灵活高效的API接口实现多渠道数据采集，结合强大的后台管理能力，提供数据展示、统计分析、报表导出等核心功能。

## 核心设计理念

- **灵活**：多种数据格式和传输方式
- **扩展**：快速响应新业务需求
- **易用**：数据录入人员和管理人员都能快速上手

## 核心能力模块

### 1. 数据采集层
- API主动推送模式
- API被动接收模式
- 支持JSON/FormData/Excel多格式

### 2. 数据管理层
- 数据清洗与校验
- 版本控制
- 数据分类管理

### 3. 业务应用层
- 数据展示与检索
- 统计分析
- 报表生成与导出

### 4. 系统服务层
- 用户认证
- 角色权限管理
- 系统配置

## 技术栈

### 前端
- Vue 3 + TypeScript
- Vite + Pinia + Vue Router
- Element Plus + Tailwind CSS

### 后端
- Go + Gin + GORM
- SQLite/MySQL/PostgreSQL

## 快速开始

### 后端启动
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 前端启动
```bash
cd frontend
npm install
npm run dev
```

### 构建
```bash
make build_all
```

## 项目结构

```
DataRegistry/
├── frontend/          # 前端项目
├── backend/           # 后端项目
├── docs/             # 文档
└── README.md
```

## License

MIT
