package main

import (
	"log"
	"time"

	"github.com/kuonji/blog/config"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	if err := repository.InitDatabase(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}

	if err := repository.AutoMigrate(); err != nil {
		logger.Fatal("Failed to auto migrate", zap.Error(err))
	}

	db := repository.GetDB()

	// ── 管理员 ──
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
		admin := model.User{
			Username: cfg.Admin.Username,
			Password: string(hash),
			Email:    cfg.Admin.Email,
			Role:     "admin",
		}
		db.Create(&admin)
		logger.Info("Admin user created", zap.String("username", cfg.Admin.Username))
	}

	// ── 分类 ──
	var catCount int64
	db.Model(&model.Category{}).Count(&catCount)
	if catCount == 0 {
		// 逐个创建以确保 ID 可用
		cats := []model.Category{
			{Name: "Go", Description: "Go 语言相关", Sort: 1},
			{Name: "Frontend", Description: "前端开发", Sort: 2},
			{Name: "AI", Description: "人工智能与机器学习", Sort: 3},
			{Name: "Notes", Description: "学习笔记与读书心得", Sort: 4},
			{Name: "TIL", Description: "Today I Learned", Sort: 5},
			{Name: "杂谈", Description: "生活杂谈与随想", Sort: 6},
		}
		for i := range cats {
			db.Create(&cats[i])
		}

		// 子分类
		db.Create(&model.Category{Name: "Gin", Description: "Gin Web 框架", ParentID: &cats[0].ID, Sort: 1})
		db.Create(&model.Category{Name: "GORM", Description: "GORM ORM", ParentID: &cats[0].ID, Sort: 2})
		db.Create(&model.Category{Name: "Vue", Description: "Vue.js 框架", ParentID: &cats[1].ID, Sort: 1})
		logger.Info("Categories seeded")
	}

	// ── 标签 ──
	var tagCount int64
	db.Model(&model.Tag{}).Count(&tagCount)
	if tagCount == 0 {
		tagNames := []string{
			"go", "gin", "gorm", "vue", "javascript", "typescript", "react",
			"ai", "llm", "python", "docker", "linux", "mysql", "redis", "git",
			"design", "tutorial", "notes", "til", "career", "life",
		}
		for _, name := range tagNames {
			db.Create(&model.Tag{Name: name})
		}
		logger.Info("Tags seeded")
	}

	// ── 系列 ──
	var seriesCount int64
	db.Model(&model.Series{}).Count(&seriesCount)
	if seriesCount == 0 {
		seriesList := []model.Series{
			{Name: "Go Web 开发实战", Description: "从零开始用 Go 构建 Web 应用的完整教程", Sort: 1},
			{Name: "Vue 3 入门到实战", Description: "Vue 3 组合式 API 的系统学习路径", Sort: 2},
			{Name: "TIL Diary", Description: "每天学一点，记录成长", Sort: 3},
		}
		for i := range seriesList {
			db.Create(&seriesList[i])
		}
		logger.Info("Series seeded")
	}

	// ── 文章 ──
	var artCount int64
	db.Model(&model.Article{}).Count(&artCount)
	if artCount == 0 {
		// 获取关联对象
		var goCat, feCat, aiCat, notesCat, tilCat, lifeCat model.Category
		db.Where("name = ? AND parent_id IS NULL", "Go").First(&goCat)
		db.Where("name = ? AND parent_id IS NULL", "Frontend").First(&feCat)
		db.Where("name = ? AND parent_id IS NULL", "AI").First(&aiCat)
		db.Where("name = ? AND parent_id IS NULL", "Notes").First(&notesCat)
		db.Where("name = ? AND parent_id IS NULL", "TIL").First(&tilCat)
		db.Where("name = ? AND parent_id IS NULL", "杂谈").First(&lifeCat)

		var goTag, ginTag, vueTag, aiTag, llmTag, notesTag, tilTag, careerTag model.Tag
		db.Where("name = ?", "go").First(&goTag)
		db.Where("name = ?", "gin").First(&ginTag)
		db.Where("name = ?", "vue").First(&vueTag)
		db.Where("name = ?", "ai").First(&aiTag)
		db.Where("name = ?", "llm").First(&llmTag)
		db.Where("name = ?", "notes").First(&notesTag)
		db.Where("name = ?", "til").First(&tilTag)
		db.Where("name = ?", "career").First(&careerTag)

		var goSeries model.Series
		db.Where("name = ?", "Go Web 开发实战").First(&goSeries)

		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		weekAgo := now.Add(-7 * 24 * time.Hour)
		monthAgo := now.Add(-30 * 24 * time.Hour)

		articles := []model.Article{
			{
				Title:       "用 Gin 搭建 RESTful API 的最佳实践",
				Content:     "# 用 Gin 搭建 RESTful API 的最佳实践\n\nGin 是 Go 语言中最流行的 Web 框架之一，以其高性能和简洁的 API 著称。本文将分享在实际项目中使用 Gin 的一些最佳实践。\n\n## 项目结构\n\n一个良好的项目结构可以让代码更易维护：\n\n```go\nbackend/\n├── cmd/\n│   └── server/\n│       └── main.go\n├── internal/\n│   ├── controller/\n│   ├── service/\n│   ├── repository/\n│   ├── model/\n│   └── middleware/\n├── config/\n└── go.mod\n```\n\n## 路由分组\n\n使用路由组来组织 API：\n\n```go\napi := r.Group(\"/api/v1\")\n{\n    public := api.Group(\"\")\n    {\n        public.GET(\"/articles\", controller.GetArticles)\n    }\n\n    admin := api.Group(\"/admin\")\n    admin.Use(middleware.JWTAuth())\n    {\n        admin.POST(\"/articles\", controller.CreateArticle)\n    }\n}\n```\n\n## 中间件\n\n合理使用中间件可以减少重复代码：\n\n- **Logger**: 记录请求日志\n- **Recovery**: 捕获 panic\n- **CORS**: 跨域处理\n- **JWT Auth**: 身份验证\n\n> 好的架构不是一开始就完美的，而是在不断重构中演进的。\n\n## 总结\n\n1. 保持项目结构清晰\n2. 合理使用中间件\n3. 统一错误处理\n4. 编写测试用例",
				Summary:     "分享使用 Gin 框架搭建 RESTful API 的项目结构、路由分组和中间件使用经验。",
				Status:      model.ArticleStatusPublished,
				PublishedAt: &now,
				AuthorID:    1,
				CategoryID:  goCat.ID,
				SeriesID:    &goSeries.ID,
				SeriesOrder: 1,
				Tags:        []model.Tag{goTag, ginTag},
			},
			{
				Title:       "GORM 高级查询技巧",
				Content:     "# GORM 高级查询技巧\n\nGORM 是 Go 语言功能强大的 ORM 库。本文整理了一些日常开发中常用的高级查询技巧。\n\n## 预加载关联\n\n```go\n// 预加载多个关联\ndb.Preload(\"Tags\").Preload(\"Category\").Find(&articles)\n\n// 条件预加载\ndb.Preload(\"Articles\", func(db *gorm.DB) *gorm.DB {\n    return db.Order(\"created_at desc\").Limit(5)\n}).Find(&categories)\n```\n\n## 复杂查询\n\n```go\n// 子查询\ndb.Where(\"category_id IN (?)\",\n    db.Model(&Category{}).Where(\"parent_id = ?\", parentID).Select(\"id\"),\n).Find(&articles)\n\n// 聚合\ndb.Model(&Article{}).Select(\"category_id, COUNT(*) as count\").\n    Group(\"category_id\").Scan(&results)\n```\n\n## 事务处理\n\n```go\ndb.Transaction(func(tx *gorm.DB) error {\n    if err := tx.Create(&article).Error; err != nil {\n        return err\n    }\n    if err := tx.Model(&article).Association(\"Tags\").Replace(tags); err != nil {\n        return err\n    }\n    return nil\n})\n```\n\n掌握这些技巧能让数据库操作更加优雅高效。",
				Summary:     "整理 GORM 的预加载、子查询、聚合和事务处理等高级用法。",
				Status:      model.ArticleStatusPublished,
				PublishedAt: &yesterday,
				AuthorID:    1,
				CategoryID:  goCat.ID,
				SeriesID:    &goSeries.ID,
				SeriesOrder: 2,
				Tags:        []model.Tag{goTag},
			},
			{
				Title:       "Vue 3 组合式 API 实战心得",
				Content:     "# Vue 3 组合式 API 实战心得\n\n从 Options API 迁移到 Composition API 后，代码组织方式有了很大变化。分享一些实际项目中的心得。\n\n## 为什么选择组合式 API\n\n- 逻辑复用更灵活（composables vs mixins）\n- 类型推导更好\n- 代码组织按功能而非选项\n\n## 常用模式\n\n```typescript\n// 数据获取 composable\nexport function useFetch<T>(url: string) {\n  const data = ref<T | null>(null)\n  const loading = ref(true)\n  const error = ref<Error | null>(null)\n\n  onMounted(async () => {\n    try {\n      data.value = await fetch(url).then(r => r.json())\n    } catch (e) {\n      error.value = e as Error\n    } finally {\n      loading.value = false\n    }\n  })\n\n  return { data, loading, error }\n}\n```\n\n## 小结\n\n组合式 API 不是银弹，但在中大型项目中优势明显。",
				Summary:     "从 Options API 迁移到 Composition API 的实战经验和常用模式分享。",
				Status:      model.ArticleStatusPublished,
				PublishedAt: &weekAgo,
				AuthorID:    1,
				CategoryID:  feCat.ID,
				Tags:        []model.Tag{vueTag},
			},
			{
				Title:       "LLM 应用开发入门：从 Prompt 到 Agent",
				Content:     "# LLM 应用开发入门\n\n大语言模型（LLM）正在改变软件开发的方式。本文从实践角度介绍如何开始 LLM 应用开发。\n\n## Prompt Engineering\n\n好的 prompt 是 LLM 应用的基础：\n\n1. **明确角色**: 告诉模型它是谁\n2. **提供上下文**: 给足背景信息\n3. **指定输出格式**: JSON、Markdown 等\n4. **给出示例**: Few-shot learning\n\n## RAG 模式\n\nRetrieval-Augmented Generation 是最常见的应用模式：\n\n```\n用户问题 → 向量检索 → 拼接上下文 → LLM 生成\n```\n\n## Agent 模式\n\n让 LLM 自主决策和执行：\n\n- 工具调用（function calling）\n- 多步推理\n- 自我反思\n\n> AI 不会取代程序员，但会用 AI 的程序员会取代不会用的。",
				Summary:     "LLM 应用开发入门指南，涵盖 Prompt Engineering、RAG 和 Agent 模式。",
				Status:      model.ArticleStatusPublished,
				PublishedAt: &monthAgo,
				AuthorID:    1,
				CategoryID:  aiCat.ID,
				Tags:        []model.Tag{aiTag, llmTag},
			},
			{
				Title:       "读书笔记：《The Pragmatic Programmer》",
				Content:     "# The Pragmatic Programmer 读书笔记\n\n重读这本经典，记录一些仍然适用的智慧。\n\n## 核心理念\n\n- **DRY**: Don't Repeat Yourself\n- **正交性**: 减少模块间的耦合\n- **曳光弹**: 先跑通最小可行路径\n- **石汤**: 先做出能用的东西，再迭代\n\n## 最喜欢的几条建议\n\n1. 你的知识是有价的资产\n2. 不要成为破窗效应的受害者\n3. 让计算机做重复的事\n4. 估算以避免意外\n\n> 经典之所以是经典，是因为它讨论的是不变的东西。",
				Summary:     "重读《The Pragmatic Programmer》的笔记，记录仍然适用的软件开发智慧。",
				Status:      model.ArticleStatusPublished,
				PublishedAt: &weekAgo,
				AuthorID:    1,
				CategoryID:  notesCat.ID,
				Tags:        []model.Tag{notesTag, careerTag},
			},
			{
				Title:       "TIL: Go 的 interface 隐式实现",
				Content:     "# TIL: Go 的 interface 隐式实现\n\n今天重新理解了 Go interface 的设计理念。\n\n## 关键点\n\nGo 的 interface 是**隐式实现**的 —— 不需要像 Java 那样写 `implements`：\n\n```go\ntype Writer interface {\n    Write(p []byte) (n int, err error)\n}\n\n// 只要实现了 Write 方法，就自动满足 Writer 接口\ntype MyWriter struct{}\nfunc (w MyWriter) Write(p []byte) (int, error) {\n    // 实现...\n    return len(p), nil\n}\n```\n\n## 为什么要这样设计？\n\n- 解耦：不需要知道对方的接口定义\n- 灵活：可以为已有类型定义新接口\n- 组合：小接口组合成大接口\n\n这就是 Go 「接受接口，返回结构」哲学的基础。",
				Summary:     "Go interface 隐式实现的设计理念和实际意义。",
				Status:      model.ArticleStatusPublished,
				PublishedAt: &yesterday,
				AuthorID:    1,
				CategoryID:  tilCat.ID,
				Tags:        []model.Tag{goTag, tilTag},
			},
			{
				Title:       "关于写博客这件事",
				Content:     "# 关于写博客这件事\n\n搭建这个博客的过程，也是一次很好的学习体验。\n\n## 为什么要写博客？\n\n1. **整理思路**: 写作是最好的思考方式\n2. **记录成长**: 回头看自己的文章很有成就感\n3. **帮助他人**: 你踩过的坑别人可能也会遇到\n4. **建立影响力**: 技术博客是最好的简历\n\n## 技术选型\n\n- 后端：Go + Gin + GORM + MySQL\n- 前端：Vue 3 + Vite + Vue Router\n- 部署：Docker + Nginx\n\n选择了前后端分离，虽然对个人博客来说有些 over-engineering，但作为学习项目很有价值。\n\n> 千里之行，始于足下。",
				Summary:     "关于为什么写博客、技术选型和搭建过程的思考。",
				Status:      model.ArticleStatusPublished,
				PublishedAt: &monthAgo,
				AuthorID:    1,
				CategoryID:  lifeCat.ID,
				Tags:        []model.Tag{careerTag},
			},
		}

		for i := range articles {
			db.Create(&articles[i])
		}
		logger.Info("Articles seeded", zap.Int("count", len(articles)))
	}

	// ── 评论 ──
	var commentCount int64
	db.Model(&model.Comment{}).Count(&commentCount)
	if commentCount == 0 {
		// 查找文章 ID
		var articles []model.Article
		db.Order("id asc").Find(&articles)
		if len(articles) >= 4 {
			c1 := model.Comment{ArticleID: articles[0].ID, Nickname: "Mira Patel", Email: "mira@example.com", Content: "这篇 Gin 最佳实践写得很清晰，项目结构部分对我帮助很大！", Status: model.CommentStatusApproved}
			db.Create(&c1)
			c2 := model.Comment{ArticleID: articles[0].ID, Nickname: "Sam Wright", Email: "sam@example.com", Content: "关于中间件那部分，能否再补充一下限流中间件的实现？", Status: model.CommentStatusApproved}
			db.Create(&c2)
			c3 := model.Comment{ArticleID: articles[2].ID, Nickname: "小明", Email: "xm@example.com", Content: "useFetch 这个 composable 的思路很好，我直接抄到项目里了。", Status: model.CommentStatusApproved}
			db.Create(&c3)
			c4 := model.Comment{ArticleID: articles[3].ID, Nickname: "Alex Chen", Email: "alex@example.com", Content: "Agent 模式那段讲得好，有没有推荐的 function calling 实践？", Status: model.CommentStatusPending}
			db.Create(&c4)

			// 回复评论
			reply := model.Comment{
				ArticleID: articles[0].ID,
				ParentID:  &c2.ID,
				Nickname:  "Kuonji",
				Email:     "admin@example.com",
				Content:   "好建议！下一篇文章我会加上限流和熔断中间件的内容。",
				Status:    model.CommentStatusApproved,
			}
			db.Create(&reply)
		}
		logger.Info("Comments seeded")
	}

	logger.Info("Seed completed successfully!")
}
