package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "glossika/docs"
	"glossika/service/internal/config"
	"go.uber.org/dig"
	"net/http"
	"time"
)

type servicePack struct {
	dig.In

	ServiceAddress config.ServiceAddress
	Handler        *gin.Engine
}

func NewServer(pack servicePack) *http.Server {
	return &http.Server{
		Addr:    pack.ServiceAddress.Glossika,
		Handler: pack.Handler,
	}
}

func NewRouterRoot(pack servicePack) *gin.RouterGroup {
	return pack.Handler.Group("glossika")
}

func NewGinEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Recovery(), cors.New(cors.Config{
		AllowMethods:    []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type"},
		MaxAge:          12 * time.Hour,
		AllowAllOrigins: true,
	}))

	return router
}

// NewBasic
// @Tags		Surprise Checker
// @version	1.0
// @produce	text/plain
// @Success	200
// @Router		/Ping [GET]
func NewBasic(pack basicPack) {
	pack.Root.GET("Ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Pong")
	})

	pack.Root.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

type basicPack struct {
	dig.In

	Root *gin.RouterGroup
}
