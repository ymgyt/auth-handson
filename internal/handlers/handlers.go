package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/server"
)

// Handlers -
type Handlers struct {
	Example *Example
	Static  *Static
}

// New -
func New(env *server.Env) *Handlers {
	return &Handlers{
		Example: &Example{Env: env},
		Static:  &Static{Env: env},
	}
}

type base struct{}

func (b base) abort(c *gin.Context, err error) {
	c.AbortWithError(http.StatusInternalServerError, err)
	fmt.Println(err)
}
