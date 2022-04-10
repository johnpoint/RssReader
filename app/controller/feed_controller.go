package controller

import (
	"RssReader/infra"
	mongoModel "RssReader/model/mongodb"
	"github.com/gin-gonic/gin"
)

type FeedListResp struct {
	ID     string
	Title  string
	Url    string
	Status int64
}

func FeedList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user = mongoModel.User{
		ID: userID.(string),
	}
	err := user.FindFeedByID(c)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	if len(user.SubFeeds) == 0 {
		returnSuccessMsg(c, "OK", []struct{}{})
		return
	}

	var feed mongoModel.Feed
	var feeds []mongoModel.Feed
	err = feed.FindByIDs(c, user.SubFeeds, feeds)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}

	returnSuccessMsg(c, "OK", feeds)
}

type SearchFeedReq struct {
	Url string `json:"url"`
}

func SearchFeed(c *gin.Context) {
	var req SearchFeedReq
	err := c.BindJSON(&req)
	if err != nil {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}

}
