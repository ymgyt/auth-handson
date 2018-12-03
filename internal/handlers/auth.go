package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/ymgyt/auth-handson/internal/app"
	"log"
	"net/http"
	"strings"
)

// Auth -
type Auth struct {
	*app.Env
	base
	auth    *auth.Client
	counter uint64
}

// User -
type User struct {
	UID string
}

type credential struct {
	Email    string
	Password string
	UID      uint64
}

type response struct {
	Op    string
	Err   string      `json:",omitempty"`
	Msg   string      `json:",omitempty"`
	Data  interface{} `json:",omitempty"`
	Token string      `json:"token,omitempty"`
	UID   string      `json:"uid,omitempty"`
}

// SignUp -
func (a *Auth) SignUp(c *gin.Context) {

	cred := readCred(c)
	ctx := context.Background()
	ur, err := a.getUser(ctx, cred)
	if err != nil {
		if auth.IsUserNotFound(err) {
			a.doCreateUser(ctx, c, cred)
			return
		}

		c.JSON(http.StatusInternalServerError, response{Err: err.Error(), Op: "firebase.auth.GetUserByEmail()"})
		return
	}

	// 既に存在するのでエラー
	c.Writer.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(c.Writer).Encode(response{
		Op:   "sign-up",
		Msg:  "user already exists",
		Data: ur,
	}); err != nil {
		panic(err)
	}
}

// SignIn -
func (a *Auth) SignIn(c *gin.Context) {
	cred := readCred(c)
	u, err := Authenticate(cred.Email, cred.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response{Op: "sign-in", Msg: "unauthenticate", Err: err.Error()})
		return
	}

	ctx := context.Background()
	token, err := a.auth.CustomToken(ctx, u.UID)
	if err != nil {
		panic(err)
	}
	c.JSON(200, response{Op: "sign-in", Msg: "successfully authenticated", Token: token, UID: u.UID})
}

func (a *Auth) doCreateUser(ctx context.Context, c *gin.Context, cred *credential) {
	created, err := a.createUser(ctx, cred)
	if err != nil {
		c.JSON(200, response{Err: err.Error(), Op: "firebase.auth.CreateUser()"})
		return
	}

	c.JSON(201, response{Op: "sign-up", Msg: "user created", Data: created})
}

func (a *Auth) createUser(ctx context.Context, cred *credential) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(cred.Email).
		EmailVerified(false).
		Password(cred.Password).
		Disabled(false)
	ur, err := a.auth.CreateUser(ctx, params)
	if err != nil {
		// need warp
		return nil, err
	}
	log.Printf("firebase user created uid: %s email: %s\n", ur.UserInfo.UID, cred.Email)
	return ur, nil
}

func (a *Auth) getUser(ctx context.Context, cred *credential) (*auth.UserRecord, error) {
	ur, err := a.auth.GetUserByEmail(ctx, cred.Email)
	if err != nil {
		return nil, err
	}
	return ur, nil
}

func readCred(c *gin.Context) *credential {
	var cred credential
	if err := c.BindJSON(&cred); err != nil {
		panic(err)
	}
	return &cred
}

// Authenticate 自前の認証のmock
func Authenticate(email, password string) (*User, error) {
	if strings.Contains(email, "@") && len(password) > 3 {
		return &User{UID: "uid-" + email}, nil
	}
	return nil, errors.New("authenticate fail")
}
