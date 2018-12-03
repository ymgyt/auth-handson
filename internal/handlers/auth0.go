package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/app"
	"golang.org/x/oauth2"
)

const (
	// TODO Remove me !
	defaultRedirectURL = "http://localhost:8123/auth0/callback"
)

// Auth0 -
type Auth0 struct {
	*app.Env
}

func (a *Auth0) oauth2Config() *oauth2.Config {
	domain := a.Env.Auth0.Domain
	return &oauth2.Config{
		ClientID:     a.Env.Auth0.ClientID,
		ClientSecret: a.Env.Auth0.ClientSecret,
		RedirectURL:  defaultRedirectURL,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}
}

// Hello -
func (a *Auth0) Hello(c *gin.Context) {
	req := c.Request
	u := req.Context().Value("user")
	spew.Dump(u)
	c.JSON(200, u)
}

// SignUp -
func (a *Auth0) SignUp(c *gin.Context) {
	cred := readCred(c)

	payload := struct {
		ClientID     string                 `json:"client_id"`
		ClientSecret string                 `json:"client_secret"`
		Email        string                 `json:"email"`
		Password     string                 `json:"password"`
		Connection   string                 `json:"connection"`
		UserMetadata map[string]interface{} `json:"user_metadata"`
	}{
		ClientID:     a.Env.Auth0.ClientID,
		ClientSecret: a.Env.Auth0.ClientSecret,
		Email:        cred.Email,
		Password:     cred.Password,
		Connection:   "auth-handson",
		UserMetadata: map[string]interface{}{
			"created": time.Now(),
		},
	}

	spew.Dump(payload)

	encoded, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}

	endpoint := a.Env.Auth0.BaseURL + "/dbconnections/signup"

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(encoded))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var got map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&got); err != nil {
		panic(err)
	}

	spew.Dump(got)

	c.JSON(resp.StatusCode, got)
}

// SignIn -
func (a *Auth0) SignIn(c *gin.Context) {
	cred := readCred(c)

	payload := struct {
		// Auth0 dashboardから確認できるapplicationとしてのcredential
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`

		GrantType string `json:"grant_type"`
		Realm     string `json:"realm"`              // userを探すconnectionを指定
		Audience  string `json:"audience,omitempty"` // 誰が得られたjwtを検証するかを指定する. Auth0側が検証する場合は指定しない

		// loginを実行したいuserのcredential
		Username string `json:"username"`
		Password string `json:"password"`
		Scope    string `json:"scope,omitempty"`
	}{
		ClientID:     a.Env.Auth0.ClientID,
		ClientSecret: a.Env.Auth0.ClientSecret,
		GrantType:    "http://auth0.com/oauth/grant-type/password-realm",
		Realm:        "auth-handson",
		//Audience:     "http://auth-handson/api",
		Username: cred.Email,
		Password: cred.Password,
	}
	spew.Dump(payload)

	encoded, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}

	endpoint := a.Env.Auth0.BaseURL + "/oauth/token"

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(encoded))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var got map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&got); err != nil {
		panic(err)
	}

	c.JSON(resp.StatusCode, got)
}

func (a *Auth0) userInfo(accessToken string) (code int, userInfo map[string]interface{}) {
	endpoint := a.Env.Auth0.BaseURL + "/userinfo"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, userInfo
	}

	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		panic(err)
	}

	return resp.StatusCode, userInfo
}

// ProtectedResource -
// jwtをvalidationするmiddlewareに保護されたendpointのsample
func (a *Auth0) ProtectedResource(c *gin.Context) {
	ctx := c.Request.Context()
	v := ctx.Value("auth0_token")
	spew.Dump("context.Value() => ", v)
	jwtToken, ok := v.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusUnauthorized, map[string]string{"msg": "failed to convert jwt token"})
		return
	}

	res := map[string]interface{}{
		"jwt":     jwtToken,
		"message": fmt.Sprintf("Hello %s", jwtToken.Claims.(jwt.MapClaims)["nickname"]),
	}

	c.JSON(200, res)
}

// OAuthLogin -
func (a *Auth0) OAuthLogin(c *gin.Context) {
	domain := a.Env.Auth0.Domain
	aud := "https://" + domain + "/userinfo"

	conf := a.oauth2Config()

	// generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.Store.Get(c.Request, "state")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session.Values["state"] = state
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	audience := oauth2.SetAuthURLParam("audience", aud)
	url := conf.AuthCodeURL(state, audience)

	c.Redirect(http.StatusSeeOther, url)
}

// Callback Auth0で認証が成功したあとにUserがredirectされてくる
func (a *Auth0) Callback(c *gin.Context) {
	domain := a.Env.Auth0.Domain

	conf := a.oauth2Config()

	state := c.Query("state")
	session, err := app.Store.Get(c.Request, "state")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if state != session.Values["state"] {
		c.AbortWithError(http.StatusInternalServerError, errors.New("invalid state parameter"))
		return
	}

	code := c.Query("code")
	token, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// fetch user info
	client := conf.Client(context.TODO(), token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session, err = app.Store.Get(c.Request, "auth-session")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/auth0/user")
}

// AppA -
func (a *Auth0) AppA(c *gin.Context) {
	v := c.Request.Context().Value("auth0_token")
	jwtToken, ok := v.(*jwt.Token)
	if !ok {
		c.String(200, "App A")
		return
	}

	res := map[string]interface{}{
		"app":     "A",
		"jwt":     jwtToken,
		"message": fmt.Sprintf("Hello %s", jwtToken.Claims.(jwt.MapClaims)["nickname"]),
	}
	c.JSON(200, res)
}

// AppB -
func (a *Auth0) AppB(c *gin.Context) {
	v := c.Request.Context().Value("auth0_token")
	jwtToken, ok := v.(*jwt.Token)
	if !ok {
		c.String(200, "App B")
		return
	}

	res := map[string]interface{}{
		"app":     "B",
		"jwt":     jwtToken,
		"message": fmt.Sprintf("Hello %s", jwtToken.Claims.(jwt.MapClaims)["nickname"]),
	}
	c.JSON(200, res)
}
