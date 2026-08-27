package controller_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kuonji/blog/internal/controller"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/testutil"
)

// ── 分类测试 ──

func TestGetCategories_Empty(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "GET", "/api/v1/categories", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"]
	if data != nil {
		list := data.([]interface{})
		if len(list) != 0 {
			t.Fatalf("expected 0 categories, got %d", len(list))
		}
	}
	// data 为 null 或空数组都算空
}

func TestGetCategories_WithArticleCount(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	testutil.CreateTestArticle(t, db, "Article 1", cat.ID, authorID, nil)

	w := testutil.DoRequest(r, "GET", "/api/v1/categories", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Fatalf("expected 1 category, got %d", len(data))
	}
	catData := data[0].(map[string]interface{})
	if catData["article_count"].(float64) != 1 {
		t.Fatalf("expected article_count 1, got %v", catData["article_count"])
	}
}

func TestGetCategory_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "GET", "/api/v1/categories/999", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateCategory_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	_, token := testutil.CreateTestAdmin(t, db)

	w := testutil.DoAuthRequest(r, "POST", "/api/v1/admin/categories",
		testutil.JSONBody(controller.CreateCategoryRequest{Name: "NewCat"}), token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCategory_WithArticles_Fails(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, token := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	testutil.CreateTestArticle(t, db, "Article", cat.ID, authorID, nil)

	w := testutil.DoAuthRequest(r, "DELETE", fmt.Sprintf("/api/v1/admin/categories/%d", cat.ID), nil, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 (has articles), got %d", w.Code)
	}
}

// ── 标签测试 ──

func TestGetTags_Empty(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "GET", "/api/v1/tags", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateTag_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	_, token := testutil.CreateTestAdmin(t, db)

	w := testutil.DoAuthRequest(r, "POST", "/api/v1/admin/tags",
		testutil.JSONBody(controller.CreateTagRequest{Name: "golang"}), token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateTag_Duplicate(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	_, token := testutil.CreateTestAdmin(t, db)
	testutil.CreateTestTag(t, db, "go")

	w := testutil.DoAuthRequest(r, "POST", "/api/v1/admin/tags",
		testutil.JSONBody(controller.CreateTagRequest{Name: "go"}), token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 (duplicate), got %d", w.Code)
	}
}

// ── 系列测试 ──

func TestGetSeries_Empty(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "GET", "/api/v1/series", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetSeriesDetail_WithArticles(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	s := testutil.CreateTestSeries(t, db, "Go Series")

	// 关联文章到系列
	art := model.Article{
		Title:      "Series Article 1",
		Content:    "content",
		CategoryID: cat.ID,
		AuthorID:   authorID,
		SeriesID:   &s.ID,
		Status:     model.ArticleStatusPublished,
	}
	db.Create(&art)

	w := testutil.DoRequest(r, "GET", fmt.Sprintf("/api/v1/series/%d", s.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	articles := data["articles"].([]interface{})
	if len(articles) != 1 {
		t.Fatalf("expected 1 article in series, got %d", len(articles))
	}
}

// ── 评论测试 ──

func TestCreateComment_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	article := testutil.CreateTestArticle(t, db, "Commented Article", cat.ID, authorID, nil)

	w := testutil.DoRequest(r, "POST", fmt.Sprintf("/api/v1/articles/%d/comments", article.ID),
		testutil.JSONBody(controller.CreateCommentRequest{
			Nickname: "TestUser",
			Email:    "test@test.com",
			Content:  "Great article!",
		}))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateComment_InvalidEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	article := testutil.CreateTestArticle(t, db, "Article", cat.ID, authorID, nil)

	w := testutil.DoRequest(r, "POST", fmt.Sprintf("/api/v1/articles/%d/comments", article.ID),
		testutil.JSONBody(controller.CreateCommentRequest{
			Nickname: "TestUser",
			Email:    "not-an-email",
			Content:  "Comment",
		}))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 (invalid email), got %d", w.Code)
	}
}

func TestGetComments_OnlyApproved(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	article := testutil.CreateTestArticle(t, db, "Article", cat.ID, authorID, nil)

	// 一条已审核 + 一条待审核
	db.Create(&model.Comment{ArticleID: article.ID, Nickname: "A", Email: "a@a.com", Content: "Approved", Status: model.CommentStatusApproved})
	db.Create(&model.Comment{ArticleID: article.ID, Nickname: "B", Email: "b@b.com", Content: "Pending", Status: model.CommentStatusPending})

	w := testutil.DoRequest(r, "GET", fmt.Sprintf("/api/v1/articles/%d/comments", article.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Fatalf("expected 1 approved comment, got %d", len(data))
	}
}

// ── 搜索测试 ──

func TestSearch_Found(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	testutil.CreateTestArticle(t, db, "Gin Framework Guide", cat.ID, authorID, nil)
	testutil.CreateTestArticle(t, db, "Vue Tutorial", cat.ID, authorID, nil)

	w := testutil.DoRequest(r, "GET", "/api/v1/search?q=Gin", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Fatalf("expected 1 result for 'Gin', got %v", data["total"])
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "GET", "/api/v1/search?q=", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 (empty query), got %d", w.Code)
	}
}

// ── 认证测试 ──

func TestLogin_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	testutil.CreateTestAdmin(t, db)

	w := testutil.DoRequest(r, "POST", "/api/v1/auth/login",
		testutil.JSONBody(controller.LoginRequest{Username: "testadmin", Password: "password"}))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["token"] == nil || data["token"] == "" {
		t.Fatal("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	testutil.CreateTestAdmin(t, db)

	w := testutil.DoRequest(r, "POST", "/api/v1/auth/login",
		testutil.JSONBody(controller.LoginRequest{Username: "testadmin", Password: "wrong"}))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// ── 统计测试 ──

func TestGetStats_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	_, token := testutil.CreateTestAdmin(t, db)

	w := testutil.DoAuthRequest(r, "GET", "/api/v1/admin/stats", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["total_articles"] == nil {
		t.Fatal("expected total_articles field")
	}
}

// ── 上传测试 ──

func TestUpload_Unauthorized(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "POST", "/api/v1/admin/upload", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUpload_NoFile(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	_, token := testutil.CreateTestAdmin(t, db)

	w := testutil.DoAuthRequest(r, "POST", "/api/v1/admin/upload", nil, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 (no file), got %d", w.Code)
	}
}
