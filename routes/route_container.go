package routes

import "github.com/gin-gonic/gin"

type RouteContainer interface {
	Register(router gin.IRouter)
}
