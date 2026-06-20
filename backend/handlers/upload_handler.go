package handlers

import (
	"io"
	"memo/config"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传图片文件"})
		return
	}
	defer file.Close()

	// 验证文件类型
	contentType := header.Header.Get("Content-Type")
	if len(contentType) < 5 || contentType[:5] != "image" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持图片文件"})
		return
	}

	// 生成文件名 (按日期分类)
	dateDir := time.Now().Format("20060102")
	filename := uuid.New().String() + filepath.Ext(header.Filename)
	uploadPath := filepath.Join(config.UploadDir, "images", dateDir, filename)

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	// 保存文件
	dst, err := os.Create(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入文件失败"})
		return
	}

	// 返回访问 URL
	c.JSON(http.StatusOK, gin.H{
		"url": "/uploads/images/" + dateDir + "/" + filename,
	})
}
