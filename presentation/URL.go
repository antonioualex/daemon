package presentation

import (
	"daemon/domain"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type URLHandler struct {
	us domain.URLService
}

type AddURLRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h URLHandler) Add(c *gin.Context) {

	URLRequest := AddURLRequest{}

	err := c.ShouldBindJSON(&URLRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err = url.ParseRequestURI(URLRequest.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url format"})
		return
	}

	err = h.us.Add(URLRequest.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url format"})
		return
	}

	h.us.Download(&domain.URL{
		URLstring: URLRequest.URL,
	},
		&domain.BatchInsights{})

	c.JSON(http.StatusOK, gin.H{"status": "URL added successfully"})
	return
}

func (h URLHandler) Retrieve(c *gin.Context) {
	sortOrder := c.Query("sort")
	switch strings.ToUpper(sortOrder) {
	case domain.ASC:
		sortOrder = domain.ASC
	case domain.DESC:
		sortOrder = domain.DESC
	case domain.BIGGEST:
		sortOrder = domain.BIGGEST
	case domain.LOWEST:
		sortOrder = domain.LOWEST
	default:
		sortOrder = ""
	}

	urls, err := h.us.Get(50, sortOrder)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, urls)
	return
}

func NewURLHandler(service domain.URLService) *URLHandler {
	return &URLHandler{us: service}
}
