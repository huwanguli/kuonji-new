package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/config"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/internal/router"
	"github.com/kuonji/blog/pkg/jwt"
	"github.com/kuonji/blog/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// SetupTestDB 初始化 SQLite 内存数据库用于测试
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// 初始化 logger（测试模式用 info console）
	logger.Init("info", "console", "stdout")

	// 初始化全局 config（JWT 中间件需要）
	config.InitTestConfig()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 与生产保持一致
		},
	})
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

	// 手动建表
	if err := db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Category{},
		&model.Tag{},
		&model.Series{},
		&model.Comment{},
		&model.VisitLog{},
	); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 注入到全局 repository
	repository.SetTestDB(db)
	return db
}

// SetupRouter 创建测试用 Gin 引擎
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.Setup("release")
}

// CreateTestAdmin 创建测试管理员并返回 token
func CreateTestAdmin(t *testing.T, db *gorm.DB) (uint, string) {
	t.Helper()
	user := model.User{
		Username: "testadmin",
		Password: "$2a$10$df4zPo2WEHypxzoXbVVjXulq6OBvS4lZLECrDTvJQqBwsgJPIRDGO", // "password"
		Email:    "test@example.com",
		Role:     "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	cfg := config.Get()
	j := jwt.New(cfg.JWT.Secret, cfg.JWT.Expire)
	token, err := j.GenerateToken(user.ID, user.Username)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	return user.ID, token
}

// CreateTestCategory 创建测试分类
func CreateTestCategory(t *testing.T, db *gorm.DB, name string) model.Category {
	t.Helper()
	cat := model.Category{Name: name, Description: name + " desc", Sort: 0}
	if err := db.Create(&cat).Error; err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	return cat
}

// CreateTestTag 创建测试标签
func CreateTestTag(t *testing.T, db *gorm.DB, name string) model.Tag {
	t.Helper()
	tag := model.Tag{Name: name}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}
	return tag
}

// CreateTestSeries 创建测试系列
func CreateTestSeries(t *testing.T, db *gorm.DB, name string) model.Series {
	t.Helper()
	s := model.Series{Name: name, Description: name + " desc"}
	if err := db.Create(&s).Error; err != nil {
		t.Fatalf("Failed to create series: %v", err)
	}
	return s
}

// CreateTestArticle 创建测试文章（已发布状态）
func CreateTestArticle(t *testing.T, db *gorm.DB, title string, catID uint, authorID uint, tags []model.Tag) model.Article {
	t.Helper()
	now := time.Now()
	article := model.Article{
		Title:       title,
		Content:     "# " + title + "\n\nSome content here.",
		Summary:     title + " summary",
		CategoryID:  catID,
		AuthorID:    authorID,
		Status:      model.ArticleStatusPublished,
		PublishedAt: &now,
	}
	if err := db.Create(&article).Error; err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}
	if len(tags) > 0 {
		db.Model(&article).Association("Tags").Replace(tags)
	}
	return article
}

// JSONBody 构建 JSON 请求体
func JSONBody(v interface{}) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// ParseResponse 解析响应 JSON
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v\nBody: %s", err, w.Body.String())
	}
	return result
}

// DoRequest 执行 HTTP 请求
func DoRequest(r http.Handler, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	if body == nil {
		body = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// DoAuthRequest 执行带认证的 HTTP 请求
func DoAuthRequest(r http.Handler, method, path string, body *bytes.Buffer, token string) *httptest.ResponseRecorder {
	if body == nil {
		body = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
