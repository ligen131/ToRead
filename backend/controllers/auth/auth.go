package auth

import (
	"errors"
	"strings"
	"time"
	"to-read/model"
	"to-read/utils/logs"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	tokenHeaderName                = "Authorization"
	accessTokenExpirationDuration  = 30 * 24 * time.Hour
	refreshTokenExpirationDuration = 30 * 24 * time.Hour
)

var jwtAccessSecretKey string
var jwtRefreshSecretKey string

type Authorization struct {
	AccessSecretKey  string `yaml:"secret-key"`
	RefreshSecretKey string `yaml:"refresh-secret-key"`
}

type Claims struct {
	UserID   uint32 `json:"user_id"`
	Role uint32 `json:"role"`
	jwt.StandardClaims
}

func InitAuthorization(a Authorization) error {
	if a.AccessSecretKey == "" {
		return errors.New("access-secret-key is empty")
	}
	jwtAccessSecretKey = a.AccessSecretKey
	if a.RefreshSecretKey == "" {
		return errors.New("refresh-secret-key is empty")
	}
	jwtRefreshSecretKey = a.RefreshSecretKey
	return nil
}

func GetJwtAccessSecretKey() string {
	return jwtAccessSecretKey
}

func GetJwtRefreshSecretKey() string {
	return jwtRefreshSecretKey
}

func GenerateAccessToken(user *model.User) (token string, expireAt time.Time, err error) {
	expireAt = time.Now().Add(accessTokenExpirationDuration)
	token, err = generateToken(user.ID, user.Role, expireAt, GetJwtAccessSecretKey())
	return token, expireAt, err
}

func GenerateRefreshToken(user *model.User) (token string, expireAt time.Time, err error) {
	expireAt = time.Now().Add(refreshTokenExpirationDuration)
	token, err = generateToken(user.ID, user.Role, expireAt, GetJwtRefreshSecretKey())
	return token, expireAt, err
}

func generateToken(userID uint32, role uint32, expireAt time.Time, secretKey string) (tokenString string, err error) {
	claims := &Claims{
		UserID: userID,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		logs.Warn("Generate token failed.", zap.Error(err))
		return "", err
	}
	return tokenString, err
}

func GetClaimsFromHeader(c echo.Context) (claims Claims, err error) {
	bearerToken := strings.Split(c.Request().Header.Get(tokenHeaderName), " ")
	if len(bearerToken) < 2 {
		return Claims{}, errors.New("invalid header")
	}
	if bearerToken[0] != "Bearer" {
		return Claims{}, errors.New("invalid header")
	}

	tokenString := bearerToken[1]
	claims = Claims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetJwtAccessSecretKey()), nil
	})
	if err != nil {
		return Claims{}, err
	}

	return claims, nil
}
