package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjj1998/url-shortener-go/internal/models"
	"github.com/kjj1998/url-shortener-go/internal/utils"
)

// GenerateShortenedUrl godoc
// @Summary generate shortened urls
// @Schemes
// @Description generate shortened urls
// @Tags example
// @Accept json
// @Produce json
// @Param longUrl body models.LongUrl true	"Add URL for shortening"
// @Success 201 {object} models.ShortenedUrl
// @Failure 400 {object} utils.HTTPError
// @Router /data/shorten [post]
func GenerateShortenedUrl(g *gin.Context) {
	var longUrlForShortening models.LongUrl
	if err := g.ShouldBindJSON(&longUrlForShortening); err != nil {
		utils.NewError(g, http.StatusBadRequest, err)
		return
	}

	if err := longUrlForShortening.Validation(); err != nil {
		utils.NewError(g, http.StatusBadRequest, err)
		return
	}

	longUrl := longUrlForShortening.LongUrl
	id := utils.GenerateUniqueId()
	shortUrl := utils.ShortenUrl(id)

	shortenedUrl := models.ShortenedUrl{
		ID:       id,
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}

	g.IndentedJSON(http.StatusCreated, shortenedUrl)
}
