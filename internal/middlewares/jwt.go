package middlewares

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTValidator -
type JWTValidator struct {
	mw *jwtmiddleware.JWTMiddleware
}

// Response -
type Response struct {
	Message string `json:"message"`
}

// Jwks -
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys -
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// NewJWTValidator -
func NewJWTValidator() gin.HandlerFunc {
	mw := jwtmiddleware.New(jwtmiddleware.Options{
		// tokenをどこから取得するかは複数設定できる
		Extractor: jwtmiddleware.FromFirst(
			jwtmiddleware.FromAuthHeader,
		),
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// 3つの検証が必要
			// https://auth0.com/docs/tokens/id-token#validate-the-claims
			// - exp
			// - iss
			// - aud

			// verify aud claim
			// Auth0側になげるので、ClientIDを設定しておく
			aud := "fVA2Tsbu0OT27v12MGur6E3uDB3pfMpx"
			// 第二引数requiredにfalseを設定しているので、tokenにaudが含まれていなければ、認証をスキップする
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}

			// verify iss claim
			// ここは必ず検証する
			iss := "https://ymgyt.auth0.com/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, true)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
		UserProperty:  "auth0_token",
		Debug:         true,
		ErrorHandler:  func(w http.ResponseWriter, r *http.Request, msg string) {}, // 勝手にレスポンスをかかせない
	})

	return func(c *gin.Context) {
		err := mw.CheckJWT(c.Writer, c.Request)
		if err != nil {
			log.Printf("middleware/jwt: %s\n", err.Error())
		}
	}
}

// ここはcacheして、認証に失敗した場合だけ再fetchして、retryするようにする
func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://ymgyt.auth0.com/.well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
