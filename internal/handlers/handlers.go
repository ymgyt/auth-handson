package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/app"
)

// Handlers -
type Handlers struct {
	Example *Example
	Static  *Static
	Auth    *Auth
	Auth0 *Auth0
}

// New -
func New(env *app.Env) *Handlers {
	return &Handlers{
		Example: &Example{Env: env},
		Static:  &Static{Env: env},
		Auth:    &Auth{Env: env, auth: env.Auth},
		Auth0: &Auth0{Env: env},
	}
}

type base struct{}

func (b base) abort(c *gin.Context, err error) {
	c.AbortWithError(http.StatusInternalServerError, err)
	fmt.Println(err)
}
