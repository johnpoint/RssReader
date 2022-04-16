package controller

import (
	"RssReader/infra"
	jwtModel "RssReader/model/jwt"
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func Login(c *gin.Context) (interface{}, error) {
	var req LoginReq
	err := c.BindJSON(&req)
	if err != nil {
		//returnErrorMsg(c, infra.ReqParseError)
		return nil, jwt.ErrMissingLoginValues
	}

	if len(req.Mail) == 0 || len(req.Password) == 0 {
		//returnErrorMsg(c, infra.ReqParseError)
		return nil, jwt.ErrMissingLoginValues
	}

	var user mongoModel.User
	err = user.FindOne(c, req.Mail, utils.Sha256(req.Password+req.Mail))
	if err != nil {
		//returnErrorMsg(c, infra.LoginError)
		return nil, jwt.ErrFailedAuthentication
	}
	return &jwtModel.User{
		UserID: user.ID,
	}, nil
}

type RegisterReq struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterReq
	err := c.BindJSON(&req)
	if err != nil {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}

	if len(req.Mail) == 0 || len(req.Password) == 0 {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}

	var user = mongoModel.User{
		Mail:     req.Mail,
		Password: utils.Sha256(req.Password + req.Mail),
		SubFeeds: []string{},
		Read:     []string{},
	}
	err = user.InsertOne(c)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	returnSuccessMsg(c, "OK", nil)
}

type ChangePasswordReq struct {
	Password string `json:"password"`
}

func ChangePassword(c *gin.Context) {
	var req ChangePasswordReq
	err := c.BindJSON(&req)
	if err != nil {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}

	userID := c.GetString("user_id")
	if len(req.Password) == 0 {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}

	var user = mongoModel.User{
		ID: userID,
	}
	err = user.FindInfoByID(c)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	err = user.UpdatePassword(c, utils.Sha256(req.Password+user.Mail))
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	returnSuccessMsg(c, "OK", nil)
}
