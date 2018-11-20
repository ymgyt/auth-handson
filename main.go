package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/ymgyt/auth-handson/internal/handlers"
	"github.com/ymgyt/auth-handson/internal/server"
)

const (
	defaultPort = "8123"
)

var (
	port string
	root string
)

func newEnv(root string) *server.Env {
	env := server.Env{AppRoot: root}

	return &env
}

func registerRoutes(env *server.Env, r *gin.Engine) {
	h := handlers.New(env)

	r.GET("/static/*filepath", h.Static.Static)
	r.GET("/js/client.js", h.Static.ClientJS)
	r.GET("/example", h.Example.Example)
}

func main() {

	if root == "" {
		fail(1, "environment variable APP_ROOT required")
	}

	env := newEnv(root)

	r := gin.Default()
	registerRoutes(env, r)

	r.Run(port)
}

func fail(code int, msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	port = ":" + port

	root = os.Getenv("APP_ROOT")
}
