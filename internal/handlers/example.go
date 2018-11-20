package handlers

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/server"
)

// Example -
type Example struct {
	*server.Env
	base
}

// Example -
func (e *Example) Example(c *gin.Context) {
	t, err := template.New("example").ParseFiles(e.AppRoot + "/internal/templates/example.tmpl")
	if err != nil {
		e.abort(c, err)
		return
	}

	var b bytes.Buffer
	if err = t.ExecuteTemplate(&b, "example.tmpl", nil); err != nil {
		e.abort(c, err)
		return
	}

	c.Writer.Write(b.Bytes())
}
