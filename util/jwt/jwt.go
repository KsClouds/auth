package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	SECRETKEY = "243223ffslsfsldfl412fdsfsdf" //私钥
	HEADER    = "Authorization"               //Cookie key
)

//解析token
func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func Get(c *gin.Context) (uint, error) {
	tokenString, err := c.Cookie(HEADER)
	if len(tokenString) == 0 {
		return 0, errors.New("已登出")
	}
	if err != nil {
		return 0, err
	}
	m, err := parseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return uint(m["uid"].(float64)), nil
}

func Set(c *gin.Context, uid uint) {
	claims := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Duration(24*7) * time.Hour).Unix(), // 过期时间，必须设置,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		fmt.Println(err)
	}
	c.SetCookie(HEADER, tokenString, 7*24*60*60, "/", "", false, true)
}

func Exp(c *gin.Context) {
	c.SetCookie(HEADER, "", 0, "/", "", false, true)
}
