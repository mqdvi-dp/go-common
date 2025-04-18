package server

import (
	"github.com/gin-gonic/gin"
)

type defaultRest struct{}

// defaultRestHandler will create an instace for default rest handler
func defaultRestHandler() *defaultRest {
	return &defaultRest{}
}

func (dr *defaultRest) Router(r *gin.RouterGroup) {}
