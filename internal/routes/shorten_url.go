package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjj1998/url-shortener-go/internal/models"
	"github.com/kjj1998/url-shortener-go/internal/repository"
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
// @Success 201 {object} models.Url
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

	shortenedUrl := models.Url{
		Id:       id,
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}

	_ = repository.Client.AddUrl(context.TODO(), shortenedUrl)

	g.IndentedJSON(http.StatusCreated, shortenedUrl)
}

// RedirectShortenedUrl godoc
// @Summary redirect shortened urls to the actual urls
// @Schemes
// @Description redirect shortened urls to the actual urls
// @Tags redirect
// @Accept json
// @Produce json
// @Param shortUrl path string false "Short URL"
// @Success 307
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /{shortUrl} [get]
func RedirectShortenedUrl(g *gin.Context) {
	shortUrl := g.Param("shortUrl")

	if shortUrl == "" {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"error": "shortUrl parameter is required"})
		return
	}

	longUrl, err := repository.Client.RetrieveUrl(context.TODO(), shortUrl)

	if err != nil {
		utils.NewError(g, http.StatusInternalServerError, err)
		return
	}
	if longUrl == "" {
		utils.NewError(g, http.StatusNotFound, errors.New("URL not found"))
		return
	}

	g.Redirect(http.StatusTemporaryRedirect, longUrl)
}
