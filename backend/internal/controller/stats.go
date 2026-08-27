package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/response"
)

type StatsResponse struct {
	TotalArticles   int64 `json:"total_articles"`
	TotalViews      int64 `json:"total_views"`
	TotalComments   int64 `json:"total_comments"`
	TotalCategories int64 `json:"total_categories"`
	TotalTags       int64 `json:"total_tags"`
	TotalSeries     int64 `json:"total_series"`
}

func GetStats(c *gin.Context) {
	db := repository.GetDB()

	var totalArticles int64
	db.Model(&model.Article{}).Where("status = ?", model.ArticleStatusPublished).Count(&totalArticles)

	var totalViews int64
	db.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)

	var totalComments int64
	db.Model(&model.Comment{}).Where("status = ?", model.CommentStatusApproved).Count(&totalComments)

	var totalCategories int64
	db.Model(&model.Category{}).Count(&totalCategories)

	var totalTags int64
	db.Model(&model.Tag{}).Count(&totalTags)

	var totalSeries int64
	db.Model(&model.Series{}).Count(&totalSeries)

	response.Success(c, StatsResponse{
		TotalArticles:   totalArticles,
		TotalViews:      totalViews,
		TotalComments:   totalComments,
		TotalCategories: totalCategories,
		TotalTags:       totalTags,
		TotalSeries:     totalSeries,
	})
}
