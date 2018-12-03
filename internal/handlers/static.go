package handlers

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/app"
)

// Static -
type Static struct {
	*app.Env
	base

	fs http.Handler
}

// ClientJS -
func (s *Static) ClientJS(c *gin.Context) {
	name := path.Base(c.Request.URL.Path)
	http.ServeFile(c.Writer, c.Request, s.AppRoot+"/client/"+name)
}

// Static -
func (s *Static) Static(c *gin.Context) {
	if s.fs == nil {
		s.fs = http.StripPrefix("/static", http.FileServer(http.Dir(s.AppRoot+"/static")))
	}

	s.fs.ServeHTTP(c.Writer, c.Request)
}
