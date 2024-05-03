package refresh_token

import (
	"api-gateway/internal/pkg/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTHandler struct {
	Sub        string
	Exp        string
	Iat        string
	Email      string
	Name       string
	Aud        []string
	Role       string
	SigningKey string
	Log        logger.Logger
	Token      string
	Timeout    int
}
type CustomClaims struct {
	*jwt.Token
	Sub  string   `json:"sub"`
	Exp  float64  `json:"exp"`
	Iat  float64  `json:"iat"`
	Aud  []string `json:"aud"`
	Role string   `json:"role"`
}

func (JWTHandler *JWTHandler) GenerateAuthJWT() (access, refresh string, err error) {
	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
		//	rtClaims     jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["sub"] = JWTHandler.Sub
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["iat"] = time.Now().Unix()
	claims["role"] = JWTHandler.Role
	claims["aud"] = JWTHandler.Aud

	access, err = accessToken.SignedString([]byte(JWTHandler.SigningKey))
	if err != nil {
		JWTHandler.Log.Error("error while generating access token", logger.Error(err))
		return
	}

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = JWTHandler.Sub
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rtClaims["iat"] = time.Now().Unix()

	refresh, err = refreshToken.SignedString([]byte(JWTHandler.SigningKey))

	if err != nil {
		JWTHandler.Log.Error("error while generating refresh token", logger.Error(err))
		return
	}

	return access, refresh, nil
}

func (jwtHandler *JWTHandler) ExtractClaims() (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)
	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtHandler.SigningKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		jwtHandler.Log.Error("invalid jwt token")
		return nil, err
	}
	return claims, nil
}

func ExtractClaim(tokenStr string, signinigKey []byte) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)
	token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return signinigKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}
	return claims, nil
}
