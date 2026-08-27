package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/config"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/internal/repository"
	"github.com/kuonji/blog/pkg/jwt"
	"github.com/kuonji/blog/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成 JWT token
	cfg := config.Get()
	j := jwt.New(cfg.JWT.Secret, cfg.JWT.Expire)
	token, err := j.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.InternalError(c, "生成令牌失败")
		return
	}

	response.Success(c, LoginResponse{
		Token: token,
		User:  user,
	})
}
