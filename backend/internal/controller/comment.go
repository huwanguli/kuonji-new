package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/response"
	"gorm.io/gorm"
)

type CreateCommentRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateCommentStatusRequest struct {
	Status int8 `json:"status" binding:"required"`
}

func GetComments(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	db := repository.GetDB()
	var comments []model.Comment

	// 获取已审核的评论，支持嵌套
	if err := db.Where("article_id = ? AND status = ?", articleID, model.CommentStatusApproved).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", model.CommentStatusApproved)
		}).
		Where("parent_id IS NULL").
		Order("created_at desc").
		Find(&comments).Error; err != nil {
		response.InternalError(c, "获取评论列表失败")
		return
	}

	response.Success(c, comments)
}

func CreateComment(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()

	// 验证文章是否存在
	var article model.Article
	if err := db.First(&article, articleID).Error; err != nil {
		response.NotFound(c, "文章不存在")
		return
	}

	// 验证父评论是否存在
	if req.ParentID != nil {
		var parentComment model.Comment
		if err := db.First(&parentComment, *req.ParentID).Error; err != nil {
			response.BadRequest(c, "父评论不存在")
			return
		}
		if parentComment.ArticleID != uint(articleID) {
			response.BadRequest(c, "父评论不属于该文章")
			return
		}
	}

	comment := model.Comment{
		ArticleID: uint(articleID),
		ParentID:  req.ParentID,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Content:   req.Content,
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Status:    model.CommentStatusPending,
	}

	if err := db.Create(&comment).Error; err != nil {
		response.InternalError(c, "创建评论失败")
		return
	}

	response.Created(c, comment)
}

func UpdateCommentStatus(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	var req UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	var comment model.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	comment.Status = req.Status

	if err := db.Save(&comment).Error; err != nil {
		response.InternalError(c, "更新评论状态失败")
		return
	}

	response.Success(c, comment)
}

func DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	db := repository.GetDB()
	var comment model.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	// 检查是否有子评论
	var childCount int64
	db.Model(&model.Comment{}).Where("parent_id = ?", commentID).Count(&childCount)
	if childCount > 0 {
		response.BadRequest(c, "该评论下有回复，无法删除")
		return
	}

	if err := db.Delete(&comment).Error; err != nil {
		response.InternalError(c, "删除评论失败")
		return
	}

	response.NoContent(c)
}
