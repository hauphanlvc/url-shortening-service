package rest

import (
	"net/http"
	"url-shortening-service/generate"

	"github.com/gin-gonic/gin"
)

type GenerateHanlder struct {
	service *generate.GenerateService
}

type UrlRequest struct {
	Url string `json:"url" binding:"required"`
}

func NewGeneateHandler(service *generate.GenerateService) *GenerateHanlder {
	return &GenerateHanlder{service: service}
}

func (gh *GenerateHanlder) Generate(c *gin.Context) {
	var urlRequest UrlRequest
	if err := c.ShouldBindJSON(&urlRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shortUrl, err := gh.service.InsertNewShortUrl(c, urlRequest.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ShortUrl": shortUrl})
}
