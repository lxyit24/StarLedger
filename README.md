# StarLedger - 星账系统

> 多租户 SaaS 财税管理平台，支持服务器租赁管理与账单结算

## 项目简介

StarLedger 是一个基于 Go + Vue 3 构建的多租户 SaaS 财务管理系统。系统采用模块化架构设计，支持按需启用/停用功能模块，为不同租户提供灵活的财税管理解决方案。

### 核心功能

- **多租户管理** - 支持多租户隔离，独立数据空间
- **RBAC 权限控制** - 基于角色的细粒度访问控制
- **服务器租赁管理** - 服务器租赁全生命周期管理（新增、编辑、续租、退租）
- **账单管理** - 账单创建、支付、取消、逾期跟踪
- **模块市场** - 按需启用/停用功能模块，动态菜单渲染
- **仪表盘** - 关键数据概览（服务器统计、到期预警、账单汇总）

## 技术栈

### 后端
- **语言**: Go 1.24+
- **框架**: Gin (HTTP)
- **ORM**: Ent (实体框架，代码生成)
- **数据库**: SQLite（默认）/ MySQL（可选）
- **认证**: JWT (golang-jwt)
- **配置**: YAML

### 前端
- **框架**: Vue 3 (Composition API + `<script setup>`)
- **构建工具**: Vite
- **UI 组件**: Element Plus
- **状态管理**: Pinia
- **HTTP 客户端**: Axios
- **路由**: Vue Router 4
- **语言**: TypeScript

## 项目结构

```
StarLedger/
├── server/                    # 后端服务
│   ├── cmd/                   # 入口文件
│   │   └── main.go
│   ├── config.yaml            # 配置文件
│   ├── ent/                   # Ent ORM
│   │   ├── schema/            # 实体 Schema 定义
│   │   ├── *.go               # 生成的 CRUD 代码
│   │   ├── hook/              # Hook 中间件
│   │   ├── migrate/           # 数据库迁移
│   │   └── runtime/           # 运行时元数据
│   ├── internal/
│   │   ├── config/            # 配置加载
│   │   ├── handler/           # HTTP 处理器
│   │   │   ├── auth/          # 认证接口
│   │   │   ├── billing/       # 账单接口
│   │   │   ├── market/        # 模块市场接口
│   │   │   ├── server/        # 服务器接口
│   │   │   ├── tenant/        # 租户接口
│   │   │   └── user/          # 用户接口
│   │   ├── middleware/        # 中间件
│   │   │   ├── cors.go        # 跨域处理
│   │   │   ├── jwt.go         # JWT 认证
│   │   │   ├── logger.go      # 请求日志
│   │   │   ├── module_access.go # 模块访问控制
│   │   │   ├── rbac.go        # 角色权限控制
│   │   │   └── tenant.go      # 租户上下文
│   │   ├── model/             # 数据模型 (DTO)
│   │   ├── module/            # 模块注册表
│   │   ├── pkg/               # 公共工具
│   │   └── service/           # 业务逻辑层
│   └── go.mod
└── web/                       # 前端应用
    ├── src/
    │   ├── api/               # API 接口封装
    │   ├── assets/            # 静态资源
    │   ├── components/        # 通用组件
    │   ├── layouts/           # 布局组件
    │   ├── router/            # 路由配置
    │   ├── stores/            # Pinia 状态管理
    │   ├── utils/             # 工具函数
    │   └── views/             # 页面视图
    │       ├── billing/       # 账单管理
    │       ├── dashboard/     # 仪表盘
    │       ├── login/         # 登录页
    │       ├── market/        # 模块市场
    │       ├── server/        # 服务器管理
    │       └── system/        # 系统管理（用户/角色/租户）
    ├── package.json
    └── vite.config.ts
```

## 快速开始

### 环境要求

- Go 1.24+
- Node.js 18+
- npm / pnpm

### 后端启动

```bash
cd server

# 安装依赖
go mod download

# 生成 Ent ORM 代码（首次或 Schema 变更后）
go generate ./ent

# 启动服务（默认端口 8080）
go run cmd/main.go
```

后端服务启动后访问 `http://localhost:8080`

### 前端启动

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器（默认端口 4080）
npm run dev
```

前端开发服务器启动后访问 `http://localhost:4080`

### 配置说明

后端配置文件 `server/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug    # debug / release

database:
  driver: sqlite # sqlite / mysql
  dsn: starledger.db

jwt:
  secret: your-secret-key
  expire: 24     # Token 有效期（小时）
```

## 功能模块

| 模块 | 标识 | 类型 | 说明 |
|------|------|------|------|
| 服务器租赁 | `server_lease` | 可选 | 服务器租赁全生命周期管理 |
| 账单管理 | `billing` | 可选 | 账单创建、支付、跟踪 |
| 用户管理 | `user` | 核心 | 用户 CRUD（不可停用） |
| 角色管理 | `role` | 核心 | 角色权限管理（不可停用） |
| 租户管理 | `tenant` | 核心 | 租户管理（不可停用） |

## API 概览

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/auth/login` | POST | 用户登录 |
| `/api/auth/register` | POST | 注册新租户 |
| `/api/auth/profile` | GET | 获取当前用户信息 |
| `/api/servers` | GET/POST | 服务器列表/创建 |
| `/api/bills` | GET/POST | 账单列表/创建 |
| `/api/market/modules` | GET | 获取模块列表 |
| `/api/users` | GET/POST | 用户列表/创建 |
| `/api/roles` | GET/POST | 角色列表/创建 |
| `/api/tenants` | GET/POST | 租户列表/创建 |

## 架构设计

### 多租户架构

- 每个租户拥有独立的数据空间
- 通过 JWT Token 中的 `tenant_id` 实现租户隔离
- 中间件自动注入租户上下文

### 模块化架构

- 模块通过注册表统一管理
- 支持运行时动态启用/停用
- 前端根据模块状态动态渲染菜单
- 后端通过中间件控制模块 API 访问

### 权限控制

- RBAC（基于角色的访问控制）
- 权限粒度到具体操作（read/write/delete）
- 支持多角色分配

## 开发指南

### 添加新功能模块

1. 在 `server/ent/schema/` 定义实体 Schema
2. 运行 `go generate ./ent` 生成代码
3. 在 `server/internal/service/` 实现业务逻辑
4. 在 `server/internal/handler/` 实现 HTTP 接口
5. 在 `server/internal/module/` 注册模块
6. 前端添加对应的 API 封装和页面组件

## License

MIT
