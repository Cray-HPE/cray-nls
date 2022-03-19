package utils

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

type ResponseError struct {
	Message string `json:"message"`
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(logger Logger) RequestHandler {
	gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return RequestHandler{Gin: engine}
}
