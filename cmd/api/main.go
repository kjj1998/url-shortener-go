package main

import (
	"encoding/binary"
	"net/http"

	"github.com/kjj1998/url-shortener-go/internal/model"
	util "github.com/kjj1998/url-shortener-go/internal/util/http"

	"github.com/gin-gonic/gin"
	"github.com/jxskiss/base62"
	"github.com/sony/sonyflake"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/kjj1998/url-shortener-go/docs"
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

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func main() {
	// id, _ := sf.NextID()
	// byteSlice := make([]byte, 8)
	// binary.BigEndian.PutUint64(byteSlice, id)
	// encoded := base62.EncodeToString(byteSlice)
	// fmt.Println(encoded)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		shortenUrl := v1.Group("/data")
		{
			shortenUrl.POST("/shorten", GenerateShortenedUrl)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	router.Run(":8080")
}

// GenerateShortenedUrl godoc
// @Summary generate shortened urls
// @Schemes
// @Description generate shortened urls
// @Tags example
// @Accept json
// @Produce json
// @Param longUrl body LongUrl true	"Add URL for shortening"
// @Success 201 {object} ShortenedUrl
// @Failure 400 {object} HTTPError
// @Router /data/shorten [post]
func GenerateShortenedUrl(g *gin.Context) {
	var longUrlForShortening model.LongUrl
	if err := g.ShouldBindJSON(&longUrlForShortening); err != nil {
		util.NewError(g, http.StatusBadRequest, err)
		return
	}

	if err := longUrlForShortening.Validation(); err != nil {
		util.NewError(g, http.StatusBadRequest, err)
		return
	}

	longUrl := longUrlForShortening.LongUrl
	id := generateUniqueId()
	shortUrl := shortenUrl(id)

	shortenedUrl := model.ShortenedUrl{
		ID:       id,
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}

	g.IndentedJSON(http.StatusCreated, shortenedUrl)
}

func generateUniqueId() uint64 {
	id, _ := sf.NextID()

	return id
}

func shortenUrl(id uint64) string {
	byteSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSlice, id)
	encoded := base62.EncodeToString(byteSlice)

	return encoded
}
