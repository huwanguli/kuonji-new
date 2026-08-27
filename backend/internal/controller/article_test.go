package controller_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kuonji/blog/internal/controller"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/testutil"
)

func TestGetArticles_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "GET", "/api/v1/articles", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	if body["code"].(float64) != 0 {
		t.Fatalf("expected code 0, got %v", body["code"])
	}
	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 0 {
		t.Fatalf("expected total 0, got %v", data["total"])
	}
	_ = db
}

func TestGetArticles_WithData(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	tag := testutil.CreateTestTag(t, db, "gin")
	testutil.CreateTestArticle(t, db, "Test Article", cat.ID, authorID, []model.Tag{tag})

	w := testutil.DoRequest(r, "GET", "/api/v1/articles?page=1&page_size=10", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Fatalf("expected total 1, got %v", data["total"])
	}
	articles := data["articles"].([]interface{})
	if len(articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(articles))
	}
	art := articles[0].(map[string]interface{})
	if art["title"] != "Test Article" {
		t.Fatalf("expected title 'Test Article', got %v", art["title"])
	}
	// 验证 category 和 tags 关联
	if art["category"] == nil {
		t.Fatal("expected category to be loaded")
	}
	tags := art["tags"].([]interface{})
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
}

func TestGetArticles_FilterByTagID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	goTag := testutil.CreateTestTag(t, db, "go")
	vueTag := testutil.CreateTestTag(t, db, "vue")
	testutil.CreateTestArticle(t, db, "Go Article", cat.ID, authorID, []model.Tag{goTag})
	testutil.CreateTestArticle(t, db, "Vue Article", cat.ID, authorID, []model.Tag{vueTag})

	w := testutil.DoRequest(r, "GET", fmt.Sprintf("/api/v1/articles?tag_id=%d", goTag.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Fatalf("expected total 1, got %v", data["total"])
	}
}

func TestGetArticles_FilterByCategoryID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	goCat := testutil.CreateTestCategory(t, db, "Go")
	feCat := testutil.CreateTestCategory(t, db, "Frontend")
	testutil.CreateTestArticle(t, db, "Go Article", goCat.ID, authorID, nil)
	testutil.CreateTestArticle(t, db, "FE Article", feCat.ID, authorID, nil)

	w := testutil.DoRequest(r, "GET", fmt.Sprintf("/api/v1/articles?category_id=%d", goCat.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Fatalf("expected total 1 for Go category, got %v", data["total"])
	}
}

func TestGetArticle_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "GET", "/api/v1/articles/999", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetArticle_Found(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, _ := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	article := testutil.CreateTestArticle(t, db, "My Article", cat.ID, authorID, nil)

	w := testutil.DoRequest(r, "GET", fmt.Sprintf("/api/v1/articles/%d", article.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["title"] != "My Article" {
		t.Fatalf("expected 'My Article', got %v", data["title"])
	}
}

func TestCreateArticle_Unauthorized(t *testing.T) {
	testutil.SetupTestDB(t)
	r := testutil.SetupRouter()

	w := testutil.DoRequest(r, "POST", "/api/v1/admin/articles",
		testutil.JSONBody(controller.CreateArticleRequest{
			Title:      "Test",
			Content:    "Content",
			CategoryID: 1,
		}))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateArticle_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	_, token := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	tag := testutil.CreateTestTag(t, db, "gin")

	w := testutil.DoAuthRequest(r, "POST", "/api/v1/admin/articles",
		testutil.JSONBody(controller.CreateArticleRequest{
			Title:      "New Article",
			Content:    "Article content here",
			Summary:    "Summary",
			CategoryID: cat.ID,
			TagIDs:     []uint{tag.ID},
			Status:     model.ArticleStatusPublished,
		}), token)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["title"] != "New Article" {
		t.Fatalf("expected 'New Article', got %v", data["title"])
	}
}

func TestUpdateArticle_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, token := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	article := testutil.CreateTestArticle(t, db, "Old Title", cat.ID, authorID, nil)

	w := testutil.DoAuthRequest(r, "PUT", fmt.Sprintf("/api/v1/admin/articles/%d", article.ID),
		testutil.JSONBody(controller.UpdateArticleRequest{
			Title:   "Updated Title",
			Content: "Updated content",
		}), token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := testutil.ParseResponse(t, w)
	data := body["data"].(map[string]interface{})
	if data["title"] != "Updated Title" {
		t.Fatalf("expected 'Updated Title', got %v", data["title"])
	}
}

func TestDeleteArticle_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, token := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	article := testutil.CreateTestArticle(t, db, "To Delete", cat.ID, authorID, nil)

	w := testutil.DoAuthRequest(r, "DELETE", fmt.Sprintf("/api/v1/admin/articles/%d", article.ID), nil, token)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	// verify gone
	w2 := testutil.DoRequest(r, "GET", fmt.Sprintf("/api/v1/articles/%d", article.ID), nil)
	if w2.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w2.Code)
	}
}

func TestUpdateArticleStatus_Publish(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := testutil.SetupRouter()
	authorID, token := testutil.CreateTestAdmin(t, db)
	cat := testutil.CreateTestCategory(t, db, "Go")
	// 创建草稿
	article := model.Article{
		Title:      "Draft",
		Content:    "Draft content",
		CategoryID: cat.ID,
		AuthorID:   authorID,
		Status:     model.ArticleStatusDraft,
	}
	db.Create(&article)

	// 发布
	w := testutil.DoAuthRequest(r, "PUT", fmt.Sprintf("/api/v1/admin/articles/%d/status", article.ID),
		testutil.JSONBody(controller.UpdateArticleStatusRequest{
			Status: model.ArticleStatusPublished,
		}), token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// 确认出现在已发布列表
	w2 := testutil.DoRequest(r, "GET", "/api/v1/articles", nil)
	body := testutil.ParseResponse(t, w2)
	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Fatalf("expected 1 published article, got %v", data["total"])
	}
}
