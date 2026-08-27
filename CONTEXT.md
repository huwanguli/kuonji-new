# Kuonji Blog — 领域模型

## 项目概述

个人学习记录和技术分享博客，前后端分离架构。

## 核心术语

### 用户 (User)
- **定义**: 博客管理员（当前仅支持单一管理员）
- **属性**: 用户名、密码哈希、邮箱、角色
- **关系**: 发布多篇文章

### 文章 (Article)
- **定义**: 博客核心内容单元，Markdown 格式存储，前端实时渲染
- **属性**: 标题、内容(Markdown)、摘要、封面图、状态、发布时间、系列顺序、访问量
- **状态**:
  - `0` draft — 草稿
  - `1` published — 已发布
  - `2` scheduled — 定时发布
- **关系**:
  - 属于一个分类 (Category)
  - 可选属于一个系列 (Series)
  - 可以有多个标签 (Tag)，多对多
  - 可以有多个评论 (Comment)

### 分类 (Category)
- **定义**: 文章的层级分类，树形结构
- **属性**: 名称、描述、父分类ID、排序
- **关系**: 可有子分类、可包含文章
- **删除规则**: 有子分类或文章时不可删除

### 标签 (Tag)
- **定义**: 文章的扁平化标记
- **属性**: 名称（唯一）
- **关系**: 与文章多对多
- **删除规则**: 有关联文章时不可删除

### 系列 (Series)
- **定义**: 相关文章的线性集合（教程/连载）
- **属性**: 名称、描述、封面图、排序
- **关系**: 包含多篇文章，文章有顺序（series_order）

### 评论 (Comment)
- **定义**: 读者对文章的反馈
- **属性**: 昵称、邮箱、内容、IP、User-Agent、状态、父评论ID
- **状态**: `0` 待审核 / `1` 已通过 / `2` 垃圾
- **关系**: 属于文章，可嵌套回复
- **特点**: 无需注册

### 访问日志 (VisitLog)
- **定义**: 文章访问记录
- **属性**: 文章ID、IP、User-Agent、时间
- **用途**: 文章阅读量统计

## 内容渲染流程

```
MySQL (Markdown 原文) → Gin API (JSON) → 前端 fetch → marked.js + highlight.js → HTML → v-html 渲染
```

- Markdown 原文存数据库，不做 HTML 预渲染
- 前端 `utils/markdown.ts` 模块顶层执行 `marked.use()`，保证 SSR 安全
- h2/h3 自动添加 `id` 属性，供 TOC 锚点使用
- 代码高亮使用 highlight.js

## 技术约束

### 前端
- **框架**: Vue 3 + Vite + Vue Router (SPA)
- **渲染**: 客户端渲染（CSR）
- **Markdown**: marked + highlight.js（前端实时渲染）
- **设计**: 编辑杂志风格（温暖奶白色调）+ 暗色模式
- **编辑器**: 左右分栏实时预览 + 图片上传（可选大小）+ 工具栏

### 后端
- **框架**: Go / Gin / GORM
- **数据库**: MySQL 8（表名单数形式，`SingularTable: true`）
- **缓存**: Redis（可选，不可用时降级）
- **认证**: JWT (HS256)
- **API**: RESTful JSON
- **文件上传**: multipart/form-data → 本地 uploads/ 目录（按年月分目录，10MB 限制）

### 测试
- **单元测试**: SQLite 内存数据库，不依赖 MySQL（30 个）
- **E2E 测试**: 脚本打运行中的服务（23 个）

### 已知问题
- GORM v1.30 AutoMigrate 与 MySQL uniqueIndex 兼容性 bug，需 `auto_migrate: false`
