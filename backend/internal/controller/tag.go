package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/response"
)

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagRequest struct {
	Name string `json:"name"`
}

func GetTags(c *gin.Context) {
	db := repository.GetDB()
	var tags []model.Tag

	if err := db.Order("name asc").Find(&tags).Error; err != nil {
		response.InternalError(c, "获取标签列表失败")
		return
	}

	response.Success(c, tags)
}

func GetTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	db := repository.GetDB()
	var tag model.Tag
	if err := db.Preload("Articles").First(&tag, id).Error; err != nil {
		response.NotFound(c, "标签不存在")
		return
	}

	response.Success(c, tag)
}

func CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()

	// 检查标签是否已存在
	var existingTag model.Tag
	if err := db.Where("name = ?", req.Name).First(&existingTag).Error; err == nil {
		response.BadRequest(c, "标签已存在")
		return
	}

	tag := model.Tag{
		Name: req.Name,
	}

	if err := db.Create(&tag).Error; err != nil {
		response.InternalError(c, "创建标签失败")
		return
	}

	response.Created(c, tag)
}

func UpdateTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	var tag model.Tag
	if err := db.First(&tag, id).Error; err != nil {
		response.NotFound(c, "标签不存在")
		return
	}

	if req.Name != "" {
		// 检查新名称是否已存在
		var existingTag model.Tag
		if err := db.Where("name = ? AND id != ?", req.Name, id).First(&existingTag).Error; err == nil {
			response.BadRequest(c, "标签名称已存在")
			return
		}
		tag.Name = req.Name
	}

	if err := db.Save(&tag).Error; err != nil {
		response.InternalError(c, "更新标签失败")
		return
	}

	response.Success(c, tag)
}

func DeleteTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	db := repository.GetDB()
	var tag model.Tag
	if err := db.First(&tag, id).Error; err != nil {
		response.NotFound(c, "标签不存在")
		return
	}

	// 检查是否有关联的文章
	var articleCount int64
	db.Model(&model.Article{}).Joins("JOIN article_tags ON article_tags.article_id = article.id").
		Where("article_tags.tag_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		response.BadRequest(c, "该标签下有文章，无法删除")
		return
	}

	if err := db.Delete(&tag).Error; err != nil {
		response.InternalError(c, "删除标签失败")
		return
	}

	response.NoContent(c)
}
