package utils

import (
	docs "github.com/Cray-HPE/cray-nls/docs"
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

	docs.SwaggerInfo.Title = "NCN Lifecycl APIs"
	docs.SwaggerInfo.Description = "<h2>This is a sample server Petstore server.</h2>\nTest"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return RequestHandler{Gin: engine}
}
