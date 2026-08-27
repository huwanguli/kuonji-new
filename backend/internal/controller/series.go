package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/response"
	"gorm.io/gorm"
)

type CreateSeriesRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Sort        int    `json:"sort"`
}

type UpdateSeriesRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Sort        *int   `json:"sort"`
}

func GetSeries(c *gin.Context) {
	db := repository.GetDB()
	var series []model.Series

	if err := db.Order("sort asc").Find(&series).Error; err != nil {
		response.InternalError(c, "获取系列列表失败")
		return
	}

	response.Success(c, series)
}

func GetSeriesDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的系列ID")
		return
	}

	db := repository.GetDB()
	var series model.Series
	if err := db.Preload("Articles", func(db *gorm.DB) *gorm.DB {
		return db.Order("series_order asc").Preload("Tags").Preload("Category")
	}).First(&series, id).Error; err != nil {
		response.NotFound(c, "系列不存在")
		return
	}

	response.Success(c, series)
}

func CreateSeries(c *gin.Context) {
	var req CreateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	series := model.Series{
		Name:        req.Name,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Sort:        req.Sort,
	}

	if err := db.Create(&series).Error; err != nil {
		response.InternalError(c, "创建系列失败")
		return
	}

	response.Created(c, series)
}

func UpdateSeries(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的系列ID")
		return
	}

	var req UpdateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	var series model.Series
	if err := db.First(&series, id).Error; err != nil {
		response.NotFound(c, "系列不存在")
		return
	}

	if req.Name != "" {
		series.Name = req.Name
	}
	if req.Description != "" {
		series.Description = req.Description
	}
	if req.CoverImage != "" {
		series.CoverImage = req.CoverImage
	}
	if req.Sort != nil {
		series.Sort = *req.Sort
	}

	if err := db.Save(&series).Error; err != nil {
		response.InternalError(c, "更新系列失败")
		return
	}

	response.Success(c, series)
}

func DeleteSeries(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的系列ID")
		return
	}

	db := repository.GetDB()
	var series model.Series
	if err := db.First(&series, id).Error; err != nil {
		response.NotFound(c, "系列不存在")
		return
	}

	// 检查是否有文章
	var articleCount int64
	db.Model(&model.Article{}).Where("series_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		response.BadRequest(c, "该系列下有文章，无法删除")
		return
	}

	if err := db.Delete(&series).Error; err != nil {
		response.InternalError(c, "删除系列失败")
		return
	}

	response.NoContent(c)
}
