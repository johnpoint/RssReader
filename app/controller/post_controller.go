package controller

import (
	"RssReader/infra"
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/errorhelper"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PostList(c *gin.Context) {
	userID := c.GetString("user_id")
	var user = &mongoModel.User{
		ID: userID,
	}
	err := user.FindFeedByID(c)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}
	limitNum := c.Param("num")
	num, err := strconv.ParseInt(limitNum, 10, 64)
	if err != nil {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}
	posts, err := new(mongoModel.Post).FindPostsByFeed(c, user.SubFeeds, num)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}
	returnSuccessMsg(c, "OK", posts)
}

func ReadPostList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user = mongoModel.User{
		ID: userID.(string),
	}
	err := user.FindReadByID(c)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	if len(user.Read) == 0 {
		returnSuccessMsg(c, "OK", []struct{}{})
		return
	}
	returnSuccessMsg(c, "OK", user.Read)
}
