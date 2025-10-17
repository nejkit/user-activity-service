package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitEngine(userHandler *UserHandler, activityHandler *ActivityHandler) *gin.Engine {
	r := gin.Default()

	//TODO: cors config
	cfg := cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(cfg))

	v1 := r.Group("/api/v1")

	users := v1.Group("/users")

	users.POST("", userHandler.HandleRegisterUser)
	users.GET("/:id", userHandler.HandleGetUserName)
	users.DELETE("/:id", userHandler.HandleRemoveUser)

	activities := v1.Group("/activities")

	activities.POST("", activityHandler.HandleRegisterActivity)
	activities.GET("/events/group", activityHandler.HandleGetActivityEvents)
	activities.GET("/events", activityHandler.HandleGetUserActivityEvents)
	activities.GET("/statistics/group", activityHandler.HandleGetActivityStatistic)
	activities.GET("/statistics", activityHandler.HandleGetUserActivityStatistic)

	return r
}
