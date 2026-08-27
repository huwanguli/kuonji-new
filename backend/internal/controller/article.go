package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/response"
	"gorm.io/gorm"
)

type CreateArticleRequest struct {
	Title       string             `json:"title" binding:"required"`
	Content     string             `json:"content" binding:"required"`
	Summary     string             `json:"summary"`
	CoverImage  string             `json:"cover_image"`
	Status      model.ArticleStatus `json:"status"`
	CategoryID  uint               `json:"category_id" binding:"required"`
	SeriesID    *uint              `json:"series_id"`
	SeriesOrder int                `json:"series_order"`
	TagIDs      []uint             `json:"tag_ids"`
}

type UpdateArticleRequest struct {
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Summary     *string `json:"summary"`
	CoverImage  *string `json:"cover_image"`
	CategoryID  uint    `json:"category_id"`
	SeriesID    *uint   `json:"series_id"`
	SeriesOrder int     `json:"series_order"`
	TagIDs      []uint  `json:"tag_ids"`
}

type UpdateArticleStatusRequest struct {
	Status      model.ArticleStatus `json:"status" binding:"required"`
	PublishedAt *string             `json:"published_at"`
}

func GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	tagID, _ := strconv.Atoi(c.Query("tag_id"))
	seriesID, _ := strconv.Atoi(c.Query("series_id"))

	offset := (page - 1) * pageSize
	db := repository.GetDB()

	var articles []model.Article
	query := db.Model(&model.Article{}).Preload("Tags").Preload("Category").Preload("Series")

	// 只查询已发布的文章
	query = query.Where("status = ?", model.ArticleStatusPublished)

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if seriesID > 0 {
		query = query.Where("series_id = ?", seriesID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = article.id").
			Where("article_tags.tag_id = ?", tagID)
	}

	var total int64
	query.Count(&total)

	if err := query.Order("published_at desc").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		response.InternalError(c, "获取文章列表失败")
		return
	}

	response.Success(c, gin.H{
		"articles":  articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	db := repository.GetDB()
	var article model.Article
	if err := db.Preload("Tags").Preload("Category").Preload("Series").First(&article, id).Error; err != nil {
		response.NotFound(c, "文章不存在")
		return
	}

	// 增加访问计数
	db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	// 记录访问日志
	visitLog := model.VisitLog{
		ArticleID: uint(id),
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	db.Create(&visitLog)

	response.Success(c, article)
}


func CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	db := repository.GetDB()

	// 验证分类是否存在
	var category model.Category
	if err := db.First(&category, req.CategoryID).Error; err != nil {
		response.BadRequest(c, "分类不存在")
		return
	}

	// 验证系列是否存在
	if req.SeriesID != nil {
		var series model.Series
		if err := db.First(&series, *req.SeriesID).Error; err != nil {
			response.BadRequest(c, "系列不存在")
			return
		}
	}

	article := model.Article{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		CoverImage:  req.CoverImage,
		AuthorID:    userID.(uint),
		CategoryID:  req.CategoryID,
		SeriesID:    req.SeriesID,
		SeriesOrder: req.SeriesOrder,
		Status:      req.Status,
	}

	// 如果直接发布，设置发布时间
	if article.Status == model.ArticleStatusPublished && article.PublishedAt == nil {
		now := time.Now()
		article.PublishedAt = &now
	}

	// 处理标签
	if len(req.TagIDs) > 0 {
		var tags []model.Tag
		db.Find(&tags, req.TagIDs)
		article.Tags = tags
	}

	if err := db.Create(&article).Error; err != nil {
		response.InternalError(c, "创建文章失败")
		return
	}

	response.Created(c, article)
}

func UpdateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	var article model.Article
	if err := db.First(&article, id).Error; err != nil {
		response.NotFound(c, "文章不存在")
		return
	}

	// 验证分类是否存在
	if req.CategoryID > 0 {
		var category model.Category
		if err := db.First(&category, req.CategoryID).Error; err != nil {
			response.BadRequest(c, "分类不存在")
			return
		}
		article.CategoryID = req.CategoryID
	}

	// 验证系列是否存在
	if req.SeriesID != nil {
		var series model.Series
		if err := db.First(&series, *req.SeriesID).Error; err != nil {
			response.BadRequest(c, "系列不存在")
			return
		}
		article.SeriesID = req.SeriesID
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Summary != nil {
		article.Summary = *req.Summary
	}
	if req.CoverImage != nil {
		article.CoverImage = *req.CoverImage
	}
	if req.SeriesID != nil {
		article.SeriesOrder = req.SeriesOrder
	}

	// 处理标签
	if req.TagIDs != nil {
		var tags []model.Tag
		db.Find(&tags, req.TagIDs)
		db.Model(&article).Association("Tags").Replace(tags)
	}

	if err := db.Save(&article).Error; err != nil {
		response.InternalError(c, "更新文章失败")
		return
	}

	response.Success(c, article)
}

func DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	db := repository.GetDB()
	var article model.Article
	if err := db.First(&article, id).Error; err != nil {
		response.NotFound(c, "文章不存在")
		return
	}

	// 删除关联的标签
	db.Model(&article).Association("Tags").Clear()

	if err := db.Delete(&article).Error; err != nil {
		response.InternalError(c, "删除文章失败")
		return
	}

	response.NoContent(c)
}

func UpdateArticleStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	var req UpdateArticleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	var article model.Article
	if err := db.First(&article, id).Error; err != nil {
		response.NotFound(c, "文章不存在")
		return
	}

	article.Status = req.Status
	if req.PublishedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.PublishedAt)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", *req.PublishedAt)
		}
		if err == nil {
			article.PublishedAt = &t
		}
	}

	if err := db.Save(&article).Error; err != nil {
		response.InternalError(c, "更新文章状态失败")
		return
	}

	response.Success(c, article)
}
