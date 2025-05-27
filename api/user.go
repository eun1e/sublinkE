package api

import (
	"log"
	"strconv"
	"sublink/models"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int
	Username string
	Nickname string
	Avatar   string
	Mobile   string
	Email    string
}

type UserAccessKey struct {
	ExpiredAt   time.Time `json:"expired_at"`
	Description string    `json:"description"`
}

// 新增用户
func UserAdd(c *gin.Context) {
	user := &models.User{
		Username: "test",
		Password: "test",
	}
	err := user.Create()
	if err != nil {
		log.Println("创建用户失败")
	}
	c.String(200, "创建用户成功")
}

// 获取用户信息
func UserMe(c *gin.Context) {
	// 获取jwt中的username
	// 返回用户信息
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	err := user.Find()
	if err != nil {
		c.JSON(400, gin.H{
			"code": "00000",
			"msg":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"avatar":   "static/avatar.gif",
			"nickname": user.Nickname,
			"userId":   user.ID,
			"username": user.Username,
			"roles":    []string{"ADMIN"},
			// "perms": []string{
			// 	"sys:menu:delete", "sys:dept:edit", "sys:dict_type:add",
			// 	"sys:dict:edit", "sys:dict:delete", "sys:dict_type:edit",
			// 	"sys:menu:add", "sys:user:add", "sys:role:edit",
			// 	"sys:dept:delete", "sys:user:password_reset", "sys:user:edit",
			// 	"sys:user:delete", "sys:dept:add", "sys:role:delete",
			// 	"sys:dict_type:delete", "sys:menu:edit", "sys:dict:add",
			// 	"sys:role:add",
			// },
		},
		"msg": "获取用户信息成功",
	})
}

// 获取所有用户
func UserPages(c *gin.Context) {
	// 获取jwt中的username
	// 返回用户信息
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	users, err := user.All()
	if err != nil {
		log.Println("获取用户信息失败")
	}
	list := []*User{}
	for i := range users {
		list = append(list, &User{
			ID:       users[i].ID,
			Username: users[i].Username,
			Nickname: users[i].Nickname,
			Avatar:   "static/avatar.gif",
		})
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"list": list,
		},
		"msg": "获取用户信息成功",
	})
}

// 更新用户信息

func UserSet(c *gin.Context) {
	NewUsername := c.Param("username")
	NewPassword := c.Param("password")
	log.Println(NewUsername, NewPassword)
	if NewUsername == "" || NewPassword == "" {
		c.JSON(400, gin.H{
			"code": "00001",
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	err := user.Set(&models.User{
		Username: NewUsername,
		Password: NewPassword,
	})
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"code": "00000",
			"msg":  err,
		})
		return
	}
	// 修改成功
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "修改成功",
	})

}

func GenerateAPIKey(c *gin.Context) {
	// 从 Header 获取用户名
	username := c.GetHeader("username")
	if username == "" {
		c.JSON(400, gin.H{"msg": "缺少用户名请求头"})
		return
	}
	user := &models.User{Username: username}
	err := user.Find()
	if err != nil {
		c.JSON(400, gin.H{"msg": "用户不存在"})
		return
	}

	var userAccessKey UserAccessKey
	if err := c.ShouldBind(&userAccessKey); err != nil {
		c.JSON(500, gin.H{"msg": "参数错误"})
		return
	}

	var accessKey models.AccessKey
	accessKey.ExpiredAt = &userAccessKey.ExpiredAt
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
		"code":      "00000",
		"accessKey": apiKey,
		"msg":       "API Key生成成功",
	})
}

func DeleteAPIKey(c *gin.Context) {

	apiKeyIDParam := c.Param("id")
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
	userIDParam := c.Param("id")
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
