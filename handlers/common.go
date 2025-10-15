package handlers

import (
	"github.com/gin-gonic/gin"
)

func InitEngine(userHandler *UserHandler, activityHandler *ActivityHandler) *gin.Engine {
	r := gin.Default()

	//TODO: cors config
	//r.Use(cors.New(cors.Config{}))

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
