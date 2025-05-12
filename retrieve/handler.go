package retrieve

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RetrieveHandler struct {
	service *RetrieveService
}

func NewRetrieveHandler(service *RetrieveService) *RetrieveHandler {
	return &RetrieveHandler{service: service}
}

func (rh *RetrieveHandler) Retrieve(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	url, err := rh.service.RetrieveShortUrl(c, shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, url.OriginalUrl)
}
