package controller

import (
	"RssReader/app/logic"
	"RssReader/infra"
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/errorhelper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
	var feeds []*mongoModel.Feed
	feeds, err = feed.FindByIDs(c, user.SubFeeds)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
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
	var feed = &mongoModel.Feed{}
	err = feed.FindByUrl(c, req.Url)
	if err != nil && err != mongo.ErrNoDocuments {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}
	if len(feed.ID) == 0 {
		feed, err = logic.Feed.GetFeed(c, req.Url)
		if err != nil {
			return
		}
	}
	returnSuccessMsg(c, "OK", feed)
}

func SubscribeFeed(c *gin.Context) {
	userID, _ := c.Get("user_id")
	feedID := c.Param("id")
	var user = &mongoModel.User{
		ID: userID.(string),
	}
	err := user.SubscribeFeed(c, feedID)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}
	returnSuccessMsg(c, "OK", nil)
}

func UnSubscribeFeed(c *gin.Context) {
	userID, _ := c.Get("user_id")
	feedID := c.Param("id")
	var user = &mongoModel.User{
		ID: userID.(string),
	}
	err := user.UnSubscribeFeed(c, feedID)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}
	returnSuccessMsg(c, "OK", nil)
}
