package main

import (
	"context"
	"fmt"
	"os"

	"firebase.google.com/go/auth"
	"github.com/ymgyt/auth-handson/internal/middlewares"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"github.com/ymgyt/auth-handson/internal/app"
	"github.com/ymgyt/auth-handson/internal/handlers"
)

const (
	defaultPort = "8123"
)

var (
	port                        string
	root                        string
	gcpServiceAccountCredential string
	auth0ClientID               string
	auth0ClientSecret           string
	auth0Domain                 string
	auth0BaseURL                string
)

func newEnv(root string) *app.Env {
	env := app.Env{
		AppRoot: root,
		Auth:    newFirebaseAppAuth(),
		Auth0: app.Auth0Env{
			Domain:       auth0Domain,
			ClientSecret: auth0ClientSecret,
			ClientID:     auth0ClientID,
			BaseURL:      auth0BaseURL,
		},
	}

	return &env
}

func registerRoutes(env *app.Env, r *gin.Engine) {
	h := handlers.New(env)

	r.GET("/static/*filepath", h.Static.Static)
	r.GET("/js/*filepath", h.Static.ClientJS)
	r.GET("/example", h.Example.Example)

	// server side auth
	r.GET("/server/auth", h.Example.ServerSideAuth)
	r.POST("/server/auth/sign-up", h.Auth.SignUp)
	r.POST("/server/auth/sign-in", h.Auth.SignIn)

	// handson
	r.GET("/handson", h.Example.HandsOn)

	// auth0
	r.POST("/auth0/sign-up", h.Auth0.SignUp)
	r.POST("/auth0/sign-in", h.Auth0.SignIn)

	protected := r.Group("/auth0/private").Use(middlewares.NewJWTValidator())
	protected.GET("/resource", h.Auth0.ProtectedResource)
	protected.GET("/app-a", h.Auth0.AppA)
	protected.GET("/app-b", h.Auth0.AppB)
}

func newFirebaseAppAuth() *auth.Client {
	opt := option.WithCredentialsJSON([]byte(gcpServiceAccountCredential))
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		DatabaseURL: "https://id-integration-handson.firebaseio.com",
		ProjectID:   "id-integration-handson",
	}, opt)
	if err != nil {
		panic(err)
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}
	return auth
}

func fail(code int, msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}

func main() {

	checkEnvironments()

	env := newEnv(root)

	r := gin.Default()
	registerRoutes(env, r)

	r.Run(port)
}

func checkEnvironments() {
	if root == "" {
		fail(1, "environment variable APP_ROOT required")
	}
	if gcpServiceAccountCredential == "" {
		fail(1, "environment variable GCP_CRED required")
	}
	if auth0ClientSecret == "" {
		fail(1, "environment variable AUTH0_CLIENT_SECRET required")
	}
	if auth0ClientID == "" {
		fail(1, "environment variable AUTH0_CLIENT_ID required")
	}
	if auth0BaseURL == "" {
		fail(1, "environment variable AUTH0_BASE_URL required")
	}
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	port = ":" + port

	root = os.Getenv("APP_ROOT")
	gcpServiceAccountCredential = os.Getenv("GCP_CRED")
	auth0ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")
	auth0ClientID = os.Getenv("AUTH0_CLIENT_ID")
	auth0Domain = os.Getenv("AUTH0_DOMAIN")
	auth0BaseURL = os.Getenv("AUTH0_BASE_URL")
}
