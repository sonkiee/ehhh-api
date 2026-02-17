package routes

import "github.com/gin-gonic/gin"

type DilemmaRoutes interface {
	Create(*gin.Context)
	List(*gin.Context)
	Get(*gin.Context)
	Vote(*gin.Context)
}

func Setup(h DilemmaRoutes) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/v1")
	{
		v1.POST("/dilemmas", h.Create)
		v1.GET("/dilemmas", h.List)
		v1.POST("/dilemmas/:id/votes", h.Vote)
		v1.GET("/dilemmas/:id", h.Get)

	}
	return r
}
