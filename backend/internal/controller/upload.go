package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/pkg/response"
)

const (
	maxUploadSize = 10 << 20 // 10MB
	uploadDir     = "uploads" // 相对于后端运行目录
)

var allowedTypes = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".webp": true, ".svg": true,
}

func UploadImage(c *gin.Context) {
	// 限制上传大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, err := c.FormFile("file")
	if err != nil {
		// MaxBytesReader 超限时返回特定错误
		if err.Error() == "http: request body too large" {
			response.BadRequest(c, "文件大小不能超过 10MB")
			return
		}
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedTypes[ext] {
		response.BadRequest(c, "不支持的文件类型，仅支持 jpg/jpeg/png/gif/webp/svg")
		return
	}

	// 检查文件大小
	if file.Size > maxUploadSize {
		response.BadRequest(c, "文件大小不能超过 10MB")
		return
	}

	// 创建上传目录（按年月分目录）
	now := time.Now()
	subDir := now.Format("2006/01")
	destDir := filepath.Join(uploadDir, subDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		response.InternalError(c, "创建上传目录失败")
		return
	}

	// 生成文件名（时间戳 + 随机后缀避免重名）
	filename := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	destPath := filepath.Join(destDir, filename)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		response.InternalError(c, "保存文件失败")
		return
	}

	// 返回可访问的 URL
	url := fmt.Sprintf("/uploads/%s/%s", subDir, filename)

	response.Success(c, gin.H{
		"url":      url,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
