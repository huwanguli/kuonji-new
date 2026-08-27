package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/response"
)

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	Sort        int    `json:"sort"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	Sort        *int   `json:"sort"`
}

func GetCategories(c *gin.Context) {
	db := repository.GetDB()
	var categories []model.Category

	// 获取所有分类（树形结构）
	if err := db.Preload("Children").Where("parent_id IS NULL").Order("sort asc").Find(&categories).Error; err != nil {
		response.InternalError(c, "获取分类列表失败")
		return
	}

	// 计算每个分类的文章数
	type CategoryWithCount struct {
		model.Category
		ArticleCount int64 `json:"article_count"`
	}
	var result []CategoryWithCount
	for _, cat := range categories {
		var count int64
		db.Model(&model.Article{}).Where("category_id = ? AND status = ?", cat.ID, model.ArticleStatusPublished).Count(&count)
		// 也计算子分类的文章数
		var childIDs []uint
		db.Model(&model.Category{}).Where("parent_id = ?", cat.ID).Pluck("id", &childIDs)
		if len(childIDs) > 0 {
			var childCount int64
			db.Model(&model.Article{}).Where("category_id IN ? AND status = ?", childIDs, model.ArticleStatusPublished).Count(&childCount)
			count += childCount
		}
		result = append(result, CategoryWithCount{Category: cat, ArticleCount: count})
	}

	response.Success(c, result)
}

func GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	db := repository.GetDB()
	var category model.Category
	if err := db.Preload("Children").Preload("Articles").First(&category, id).Error; err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	response.Success(c, category)
}

func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()

	// 验证父分类是否存在
	if req.ParentID != nil {
		var parent model.Category
		if err := db.First(&parent, *req.ParentID).Error; err != nil {
			response.BadRequest(c, "父分类不存在")
			return
		}
	}

	category := model.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
	}

	if err := db.Create(&category).Error; err != nil {
		response.InternalError(c, "创建分类失败")
		return
	}

	response.Created(c, category)
}

func UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	var category model.Category
	if err := db.First(&category, id).Error; err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	// 验证父分类是否存在（不能设置自己为父分类）
	if req.ParentID != nil {
		if *req.ParentID == uint(id) {
			response.BadRequest(c, "不能设置自己为父分类")
			return
		}
		var parent model.Category
		if err := db.First(&parent, *req.ParentID).Error; err != nil {
			response.BadRequest(c, "父分类不存在")
			return
		}
		category.ParentID = req.ParentID
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Sort != nil {
		category.Sort = *req.Sort
	}

	if err := db.Save(&category).Error; err != nil {
		response.InternalError(c, "更新分类失败")
		return
	}

	response.Success(c, category)
}

func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	db := repository.GetDB()
	var category model.Category
	if err := db.First(&category, id).Error; err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	// 检查是否有子分类
	var childCount int64
	db.Model(&model.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		response.BadRequest(c, "该分类下有子分类，无法删除")
		return
	}

	// 检查是否有文章
	var articleCount int64
	db.Model(&model.Article{}).Where("category_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		response.BadRequest(c, "该分类下有文章，无法删除")
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		response.InternalError(c, "删除分类失败")
		return
	}

	response.NoContent(c)
}
