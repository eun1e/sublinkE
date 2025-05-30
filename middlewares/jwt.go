package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
	"sublink/models"
	"sublink/utils"
)

var Secret = []byte("sublink") // 秘钥

// JwtClaims jwt声明
type JwtClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthToken 验证token中间件
func AuthToken(c *gin.Context) {
	// 定义白名单
	list := []string{"/static", "/api/v1/auth/login", "/api/v1/auth/captcha", "/c/", "/api/v1/version"}
	// 如果是首页直接跳过
	if c.Request.URL.Path == "/" {
		c.Next()
		return
	}
	// 如果是白名单直接跳过
	for _, v := range list {
		if strings.HasPrefix(c.Request.URL.Path, v) {
			c.Next()
			return
		}
	}

	// 检查api key
	accessKey := c.GetHeader("X-API-Key")

	if accessKey != "" {
		username, bool, err := validApiKey(accessKey)
		if err != nil || !bool {
			c.JSON(400, gin.H{"msg": "无效的API Key"})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
		return
	}

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(400, gin.H{"msg": "请求未携带token"})
		c.Abort()
		return
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		c.JSON(400, gin.H{"msg": "token格式错误"})
		c.Abort()
		return
	}
	// 去掉Bearer前缀
	token = strings.Replace(token, "Bearer ", "", -1)
	mc, err := ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  err.Error(),
		})
		c.Abort()
		return
	}
	c.Set("username", mc.Username)
	c.Next()
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*JwtClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func validApiKey(apiKey string) (string, bool, error) {

	// 快速格式验证
	parts := strings.Split(apiKey, "_")
	if len(parts) != 3 {
		return "", false, fmt.Errorf("API Key格式错误")
	}

	encryptionKey := os.Getenv("API_ENCRYPTION_KEY")
	if encryptionKey == "" {
		return "", false, fmt.Errorf("未设置API_ENCRYPTION_KEY环境变量")
	}

	// 解密用户ID
	userID, err := utils.DecryptUserIDCompact(parts[1], []byte(encryptionKey))
	if err != nil {
		return "", false, fmt.Errorf("解密用户ID失败: %w", err)
	}

	// 数据库查询
	keys, err := models.FindValidAccessKeys(userID)
	if err != nil {
		return "", false, fmt.Errorf("查询Access Key失败: %w", err)
	}

	// bcrypt验证
	for _, key := range keys {
		if key.VerifyKey(apiKey) {

			return key.Username, true, nil
		}
	}

	return "", false, fmt.Errorf("无效的API Key")
}
