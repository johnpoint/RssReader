package apis

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"rssreader/src/model"
	"time"
)

func Login(c echo.Context) error {
	conf := model.Config{}
	err := conf.Load()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	salt := conf.Salt
	u := model.User{}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	user := u
	err = u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "account or password incorrect"})
	}
	data := []byte(user.Mail + salt + user.Password)
	has := md5.Sum(data)
	md5Password := fmt.Sprintf("%x", has)
	if !u.VerPassword(md5Password) {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "account or password incorrect"})
	}
	claims := &model.JwtCustomClaims{
		Mail: u.Mail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(salt))
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: t})
}

func Register(c echo.Context) error {
	conf := model.Config{}
	err := conf.Load()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	salt := conf.Salt
	u := model.User{}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	if u.Mail == "" || u.Password == "" {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "Password or mail can not be blank"})
	}
	data := []byte(u.Mail + salt + u.Password)
	has := md5.Sum(data)
	md5Password := fmt.Sprintf("%x", has)
	u.Password = md5Password
	err = u.New()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func CheckAuth(c echo.Context) (model.User, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	u := model.User{Mail: claims.Mail}
	err := u.Get()
	if err != nil {
		return model.User{}, errors.New("user does not exist")
	}
	return u, nil
}

func JwtError(error) error {
	return echo.NewHTTPError(http.StatusUnauthorized, model.Response{Code: 0, Message: "invalid or expired jwt"})
}
