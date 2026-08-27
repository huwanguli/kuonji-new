# Kuonji Blog

个人学习记录与技术分享博客，前后端分离架构。

## 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.21+ / Gin / GORM / MySQL 8 |
| 前端 | Vue 3 / Vite / Vue Router (SPA) |
| 认证 | JWT (HS256) |
| 日志 | zap (结构化日志) |
| 配置 | YAML + 环境变量 |
| 测试 | go test (SQLite 内存DB) / E2E 脚本 |

## 功能特性

### 前台
- 文章列表（分页、按分类/标签/系列过滤）
- 文章详情（Markdown 渲染 + 代码高亮 + 右侧目录 TOC + 系列文章目录）
- 分类浏览（树形结构，含文章数统计）
- 标签浏览
- 系列浏览（线性排序 + 上下章导航）
- 全文搜索
- 评论（无需注册，昵称+邮箱，嵌套回复，审核机制）
- 暗色模式（跟随系统或手动切换，localStorage 持久化）
- 首页 Hero 区域（渐变横幅 + 简介 + 实时统计）
- 阅读信息（字数 + 预估时间，中文显示）
- 响应式设计（移动端适配）

### 后台管理 (`/admin`)
- JWT 登录认证
- 仪表盘（6 项数据统计）
- 文章管理（创建/编辑/删除/发布草稿）
- Markdown 编辑器（左右分栏实时预览 + 滚动同步）
- 编辑器工具栏（H1-H3/粗体/斜体/删除线/代码/代码块/链接/引用/列表/分割线/图片上传）
- 图片上传（multipart → 本地存储，按年月分目录，10MB 限制，可选大小插入）
- 封面图管理（上传或 URL）
- 分类管理（添加/删除，树形结构）
- 标签管理（添加/删除）
- 系列管理（添加/删除）
- 评论管理（审核通过/标记垃圾/删除）

## 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 1. 启动后端

```bash
cd backend

# 编辑配置（数据库连接等）
vim config/config.yaml

# 安装依赖
go mod download

# 初始化数据库和种子数据
go run ./cmd/seed

# 启动服务
go run ./cmd/server
# → http://localhost:8080
```

> **注意**：`config.yaml` 中 `database.auto_migrate` 默认为 `false`。首次部署时通过 `go run ./cmd/seed` 初始化表结构。如需开启自动迁移，设为 `true`，但 GORM v1.30 在 MySQL 上有已知的索引命名兼容问题。

### 2. 启动前端

```bash
cd frontend

npm install
npm run dev
# → http://localhost:3000 (API 自动代理到 :8080)
```

### 3. 登录后台

访问 `http://localhost:3000/admin/login`

- 用户名: `admin`
- 密码: `admin123`

### 4. 运行测试

```bash
# 单元测试（SQLite 内存数据库，无需 MySQL）
cd backend
go test -v ./internal/controller/

# E2E 测试（需要后端正在运行）
pwsh -File e2e_test.ps1     # Windows
bash e2e_test.sh             # Linux/Mac
```

## 项目结构

```
kuonji/
├── backend/                        # Go 后端
│   ├── cmd/
│   │   ├── server/main.go          # 服务入口
│   │   └── seed/main.go            # 数据初始化（7 篇示例文章）
│   ├── config/
│   │   ├── config.go               # 配置加载
│   │   └── config.yaml             # 配置文件
│   ├── internal/
│   │   ├── controller/             # HTTP 控制器 + 单元测试
│   │   │   ├── article.go          # 文章 CRUD
│   │   │   ├── category.go         # 分类 CRUD
│   │   │   ├── tag.go              # 标签 CRUD
│   │   │   ├── series.go           # 系列 CRUD
│   │   │   ├── comment.go          # 评论管理
│   │   │   ├── upload.go           # 图片上传
│   │   │   ├── search.go           # 搜索
│   │   │   ├── auth.go             # 登录认证
│   │   │   ├── stats.go            # 数据统计
│   │   │   ├── article_test.go     # 文章单元测试 (11)
│   │   │   └── public_test.go      # 公开 API + 管理 API 单元测试 (19)
│   │   ├── middleware/              # JWT 认证、日志中间件
│   │   ├── model/                  # 数据模型（GORM）
│   │   ├── repository/             # 数据库连接
│   │   ├── router/                 # 路由定义 + CORS + 静态文件
│   │   └── testutil/               # 测试基础设施（SQLite 内存 DB）
│   ├── pkg/
│   │   ├── jwt/                    # JWT 工具
│   │   ├── logger/                 # zap 日志
│   │   └── response/               # 统一响应格式
│   ├── uploads/                    # 上传的图片（gitignore）
│   ├── e2e_test.ps1                # E2E 测试 (Windows)
│   └── e2e_test.sh                 # E2E 测试 (Linux)
│
├── frontend/                       # Vue 3 前端
│   ├── src/
│   │   ├── assets/css/main.css     # 设计系统（亮色 + 暗色模式）
│   │   ├── components/
│   │   │   ├── ArticleCard.vue     # 文章卡片
│   │   │   ├── ArticleTOC.vue      # 文章目录（滚动高亮 + 点击跳转）
│   │   │   ├── ArticleComments.vue # 评论区（表单 + 列表）
│   │   │   ├── CommentThread.vue   # 评论线程（嵌套回复）
│   │   │   ├── CategoryNode.vue    # 分类树节点
│   │   │   ├── PageHero.vue        # 首页 Hero 区域
│   │   │   └── SidebarLayout.vue   # 侧边栏 + 内容区双栏布局
│   │   ├── pages/                  # 前台页面
│   │   ├── pages/admin/            # 后台管理页面
│   │   │   ├── AdminLayout.vue     # 后台布局（侧边导航）
│   │   │   ├── LoginPage.vue       # 登录页
│   │   │   ├── DashboardPage.vue   # 仪表盘
│   │   │   ├── ArticlesManage.vue  # 文章管理列表
│   │   │   ├── ArticleEditor.vue   # Markdown 编辑器 + 实时预览
│   │   │   ├── CategoriesManage.vue
│   │   │   ├── TagsManage.vue
│   │   │   ├── SeriesManage.vue
│   │   │   └── CommentsManage.vue
│   │   ├── utils/
│   │   │   ├── api.ts              # API 封装（含图片上传）
│   │   │   ├── auth.ts             # JWT 认证状态管理
│   │   │   ├── colors.ts           # 分类色点映射
│   │   │   ├── format.ts           # 日期/字数/阅读时间/相对时间
│   │   │   ├── markdown.ts         # Markdown 渲染（marked + hljs）
│   │   │   └── theme.ts            # 暗色模式管理
│   │   ├── routes.ts               # 路由表
│   │   ├── App.vue                 # 根组件（导航 + 暗色切换）
│   │   └── main.ts                 # 入口
│   ├── index.html
│   └── vite.config.ts              # Vite 配置 + API/uploads 代理
│
├── docker/                         # Docker 配置（待更新）
├── README.md
└── CONTEXT.md                      # 领域模型文档
```

## API 概览

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/articles` | 文章列表（`page` `page_size` `category_id` `tag_id` `series_id`） |
| GET | `/api/v1/articles/:id` | 文章详情 |
| GET | `/api/v1/categories` | 分类列表（树形，含文章数） |
| GET | `/api/v1/categories/:id` | 分类详情 |
| GET | `/api/v1/tags` | 标签列表 |
| GET | `/api/v1/tags/:id` | 标签详情 |
| GET | `/api/v1/series` | 系列列表 |
| GET | `/api/v1/series/:id` | 系列详情（含文章列表） |
| GET | `/api/v1/articles/:id/comments` | 评论列表（仅已审核） |
| POST | `/api/v1/articles/:id/comments` | 发表评论 |
| GET | `/api/v1/search?q=keyword` | 全文搜索 |
| POST | `/api/v1/auth/login` | 管理员登录 |
| GET | `/uploads/...` | 静态图片文件 |

### 管理接口（需 JWT）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/upload` | 图片上传（multipart/form-data，10MB 限制） |
| GET | `/api/v1/admin/stats` | 数据统计 |
| POST | `/api/v1/admin/articles` | 创建文章 |
| PUT | `/api/v1/admin/articles/:id` | 更新文章 |
| DELETE | `/api/v1/admin/articles/:id` | 删除文章 |
| PUT | `/api/v1/admin/articles/:id/status` | 更新文章状态 |
| POST | `/api/v1/admin/categories` | 创建分类 |
| PUT | `/api/v1/admin/categories/:id` | 更新分类 |
| DELETE | `/api/v1/admin/categories/:id` | 删除分类 |
| POST | `/api/v1/admin/tags` | 创建标签 |
| PUT | `/api/v1/admin/tags/:id` | 更新标签 |
| DELETE | `/api/v1/admin/tags/:id` | 删除标签 |
| POST | `/api/v1/admin/series` | 创建系列 |
| PUT | `/api/v1/admin/series/:id` | 更新系列 |
| DELETE | `/api/v1/admin/series/:id` | 删除系列 |
| PUT | `/api/v1/admin/comments/:id/status` | 更新评论状态 |
| DELETE | `/api/v1/admin/comments/:id` | 删除评论 |

## 已知问题

- **GORM v1.30 AutoMigrate**：与 MySQL 的 `uniqueIndex` 存在兼容性 bug（会尝试 `DROP CONSTRAINT` 但 MySQL 将其解释为 `DROP FOREIGN KEY`）。解决方案：`database.auto_migrate: false`，通过 `seed` 命令初始化表结构。
