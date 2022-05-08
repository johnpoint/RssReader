package controller

import (
	"RssReader/infra"
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/errorhelper"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PostListReq struct {
	ID          string `json:"_id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	FID         string `json:"fid" bson:"fid"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Published   int64  `json:"published" bson:"published"`
	FTitle      string `json:"fTitle"`
}

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
	feeds, err := new(mongoModel.Feed).FindByIDs(c, user.SubFeeds)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}

	var feedMap = make(map[string]string)
	for _, v := range feeds {
		feedMap[v.ID] = v.Title
	}

	var postsRes []*PostListReq

	for _, v := range posts {
		fTitle, _ := feedMap[v.FID]
		postsRes = append(postsRes, &PostListReq{
			ID:        v.ID,
			Title:     v.Title,
			FID:       v.FID,
			Url:       v.Url,
			Published: v.Published,
			FTitle:    fTitle,
		})
	}
	returnSuccessMsg(c, "OK", postsRes)
}

func ReadPostList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var read = &mongoModel.Read{
		UId: userID.(string),
	}
	readList, err := read.FindReadListByUserId(c)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}
	returnSuccessMsg(c, "OK", readList)
}

func GetPostContent(c *gin.Context) {
	postID := c.Param("id")
	var post = mongoModel.Post{
		ID: postID,
	}
	err := post.FindPostByID(c)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	returnSuccessMsg(c, "OK", post)
}

func PostAsRead(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("user_id")

	var read = &mongoModel.Read{
		UId: userID,
		PId: postID,
	}
	err := read.MarkAsRead(c, []*mongoModel.Read{read})
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	returnSuccessMsg(c, "OK", nil)
}

func PostAsUnRead(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("user_id")
	var read = &mongoModel.Read{
		UId: userID,
		PId: postID,
	}
	err := read.MarkAsUnRead(c, []*mongoModel.Read{read})
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	returnSuccessMsg(c, "OK", nil)
}
