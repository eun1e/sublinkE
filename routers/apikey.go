package routers

import (
	"sublink/api"

	"github.com/gin-gonic/gin"
)

func APIKey(r *gin.Engine) {
	userGroup := r.Group("/api/v1/apikey")
	{
		userGroup.POST("/add", api.GenerateAPIKey)
		userGroup.DELETE("/delete", api.DeleteAPIKey)
		userGroup.GET("/get", api.GetAPIKey)
	}
}
