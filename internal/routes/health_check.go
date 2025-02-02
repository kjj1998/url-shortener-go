package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjj1998/url-shortener-go/internal/models"
)

// HealthCheck godoc
// @Summary API healthcheck
// @Schemes
// @Description Check API health
// @Produce json
// @Success 200
// @Router /health [get]
func HealthCheck(g *gin.Context) {

	health := models.HealthCheck{Health: "OK"}

	g.IndentedJSON(http.StatusOK, health)
}
