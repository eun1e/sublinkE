package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sublink/models"
	"time"
)

type UserAccessKey struct {
	UserName    string `binding:"required"`
	ExpiredAt   *time.Time
	Description string
}

func GenerateAPIKey(c *gin.Context) {
	var userAccessKey UserAccessKey
	if err := c.ShouldBindJSON(&userAccessKey); err != nil {
		c.JSON(500, gin.H{"msg": "参数错误"})
		return
	}
	user := &models.User{Username: userAccessKey.UserName}
	err := user.Find()
	if err != nil {
		c.JSON(400, gin.H{"msg": "用户不存在"})
		return
	}

	var accessKey models.AccessKey
	accessKey.ExpiredAt = userAccessKey.ExpiredAt
	accessKey.Description = userAccessKey.Description
	accessKey.UserID = user.ID
	accessKey.CreatedAt = time.Now()
	accessKey.Username = user.Username

	apiKey, err := accessKey.GenerateAPIKey()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"msg": "生成API Key失败"})
		return
	}
	err = accessKey.Generate()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"msg": "生成API Key失败"})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": map[string]string{
			"accessKey": apiKey,
		},
		"msg": "API Key生成成功",
	})
}

func DeleteAPIKey(c *gin.Context) {

	apiKeyIDParam := c.Param("apiKeyId")
	if apiKeyIDParam == "" {
		c.JSON(400, gin.H{"msg": "缺少API Key ID"})
		return
	}

	var accessKey models.AccessKey
	apiKeyID, err := strconv.Atoi(apiKeyIDParam)
	if err != nil {
		c.JSON(500, gin.H{"msg": "删除API Key失败"})
		return
	}
	accessKey.ID = apiKeyID
	err = accessKey.Delete()
	if err != nil {
		c.JSON(500, gin.H{"msg": "删除API Key失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "删除API Key成功",
	})

}

func GetAPIKey(c *gin.Context) {
	userIDParam := c.Param("userId")
	if userIDParam == "" {
		c.JSON(400, gin.H{"msg": "缺少User ID"})
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(500, gin.H{"msg": "删除API Key失败"})
		return
	}
	apiKeys, err := models.FindValidAccessKeys(userID)
	if err != nil {
		c.JSON(500, gin.H{"msg": "查询API Key失败"})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": apiKeys,
		"msg":  "查询API Key成功",
	})
}
