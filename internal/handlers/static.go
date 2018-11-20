package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/server"
)

// Static -
type Static struct {
	*server.Env
	base

	fs http.Handler
}

// ClientJS -
func (s *Static) ClientJS(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, s.AppRoot+"/client/client.js")
}

// Static -
func (s *Static) Static(c *gin.Context) {
	if s.fs == nil {
		s.fs = http.StripPrefix("/static", http.FileServer(http.Dir(s.AppRoot+"/static")))
	}

	s.fs.ServeHTTP(c.Writer, c.Request)
}
