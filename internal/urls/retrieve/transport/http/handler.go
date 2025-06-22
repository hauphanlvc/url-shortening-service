package http

import (
	"errors"
	"net/http"
	"url-shortening-service/internal/repository"
	retrieve "url-shortening-service/internal/urls/retrieve/service"

	"github.com/gin-gonic/gin"
)

type RetrieveHandler struct {
	service *retrieve.RetrieveService
}

func NewRetrieveHandler(service *retrieve.RetrieveService) *RetrieveHandler {
	return &RetrieveHandler{service: service}
}

func (rh *RetrieveHandler) Retrieve(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	url, err := rh.service.RetrieveShortUrl(c, shortUrl)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound): // if your service exposes this
			c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.Redirect(http.StatusFound, url)
}
