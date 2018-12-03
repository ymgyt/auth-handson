package handlers

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/app"
)

// Example -
type Example struct {
	*app.Env
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

// ServerSideAuth -
func (e *Example) ServerSideAuth(c *gin.Context) {
	t, err := template.New("server_side_auth").ParseFiles(e.AppRoot + "/internal/templates/server_side_auth.tmpl")
	if err != nil {
		e.abort(c, err)
		return
	}

	if err = t.ExecuteTemplate(c.Writer, "server_side_auth.tmpl", nil); err != nil {
		e.abort(c, err)
		return
	}
}

func (e *Example) HandsOn(c *gin.Context) {
	t, err := template.New("handson").ParseFiles(e.AppRoot + "/internal/templates/handson.tmpl")
	if err != nil {
		e.abort(c, err)
		return
	}

	if err = t.ExecuteTemplate(c.Writer, "handson.tmpl", nil); err != nil {
		e.abort(c, err)
		return
	}
}
