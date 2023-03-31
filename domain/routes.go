package domain

import (
	"github.com/gin-gonic/gin"
)

type RouteDef struct {
	Method      string
	HandlerFunc gin.HandlerFunc
}
