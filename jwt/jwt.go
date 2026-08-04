package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"stock-tracker/config"
	"stock-tracker/constant"
	"stock-tracker/response"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserLogin struct {
	ExpiredDt     time.Time `json:"expiredDt"`
	UserID        string    `json:"userId"`
	PassVersion   int       `json:"passVersion"`
	CompanyID     string    `json:"companyId"`
	Role          string    `json:"role"`
	UsercompanyID string    `json:"usercompanyId"`
}

func GetUserLoginInfo(c echo.Context) (UserLogin, error) {
	if u, ok := c.Get(constant.TokenUserContext).(UserLogin); ok {
		return u, nil
	} else {
		return UserLogin{}, response.ErrorForce(http.StatusUnauthorized, response.ErrorUnauthorized).SendJSON(c)
	}
}

func CreateToken(userLogin UserLogin) (string, error) {
	var err error

	now := time.Now()
	expiredUnix := userLogin.ExpiredDt.Unix()
	subject := fmt.Sprintf("%d$$%s$$%d$$%s$$%s$$%s", expiredUnix, userLogin.UserID, userLogin.PassVersion, userLogin.CompanyID, userLogin.Role, userLogin.UsercompanyID)
	jwtKey := []byte(config.JwtSecretKey)
	claims := jwt.MapClaims{
		"iss": "Services",
		"sub": subject,
		"iat": now.Unix(),
		"exp": expiredUnix,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ExtractClaims(header string) (UserLogin, error) {
	var err error
	var userlogin UserLogin

	if header == "" {
		err = errors.New("unauthorized.")
		return userlogin, err
	}

	token := header[(len(constant.BearerSchema) + 1):]
	claims, err := parseToken(token)
	if err != nil {
		return userlogin, err
	}

	content := claims["sub"].(string)
	contentData := strings.Split(content, "$$")
	if len(contentData) != constant.TokenContentLen {
		err = errors.New("unauthorized!")
		return userlogin, err
	}

	expiredUnix, err := strconv.ParseInt(contentData[0], 10, 64)
	if err != nil {
		return userlogin, err
	}

	passVersion, err := strconv.Atoi(contentData[2])
	if err != nil {
		return userlogin, err
	}

	expiredAt := time.Unix(expiredUnix, 0)
	now := time.Now()
	if now.After(expiredAt) {
		err = errors.New("token expired")
		return userlogin, err
	}
	userlogin = UserLogin{
		ExpiredDt:     expiredAt,
		UserID:        contentData[1],
		PassVersion:   passVersion,
		CompanyID:     contentData[3],
		Role:          contentData[4],
		UsercompanyID: contentData[5],
	}

	return userlogin, err
}

func parseToken(token string) (jwt.MapClaims, error) {

	jwtKey := []byte(config.JwtSecretKey)
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
