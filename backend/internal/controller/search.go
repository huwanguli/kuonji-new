package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/response"
)

func Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.BadRequest(c, "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	db := repository.GetDB()
	var articles []model.Article

	// 简单搜索：在标题和内容中搜索
	query := db.Model(&model.Article{}).Preload("Tags").Preload("Category").
		Where("status = ?", model.ArticleStatusPublished).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	var total int64
	query.Count(&total)

	if err := query.Order("published_at desc").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		response.InternalError(c, "搜索失败")
		return
	}

	response.Success(c, gin.H{
		"articles":  articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"keyword":   keyword,
	})
}
