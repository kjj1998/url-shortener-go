package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/kjj1998/url-shortener-go/docs"
	"github.com/kjj1998/url-shortener-go/internal/routes"
)

//	@title			URL shortening API
//	@version		1.0
//	@description	This is a URL shortener service.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:80"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		shortenUrl := v1.Group("/data")
		{
			shortenUrl.POST("/shorten", routes.GenerateShortenedUrl)
		}
		v1.GET("/:shortUrl", routes.RedirectShortenedUrl)
		v1.GET("/health", routes.HealthCheck)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:80/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	router.Run(":80")
}
