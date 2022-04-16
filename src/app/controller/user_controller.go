package controller

import (
	"RssReader/app/logic"
	"RssReader/infra"
	mongoModel "RssReader/model/mongodb"
	"RssReader/pkg/errorhelper"
	"RssReader/pkg/log"
	"context"
	"encoding/xml"
	"github.com/gilliek/go-opml/opml"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func ImportOPML(c *gin.Context) {
	userID := c.GetString("user_id")
	file, err := c.FormFile("opml")
	if err != nil {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}
	src, err := file.Open()
	if err != nil {
		returnErrorMsg(c, infra.ReqParseError)
		return
	}
	defer src.Close()
	all, err := io.ReadAll(src)
	if err != nil {
		return
	}
	doc, err := opml.NewOPML(all)
	var feedUrls []string
	allItem := GetOutlinesItem(doc.Body.Outlines)
	for _, i := range allItem {
		feedUrls = append(feedUrls, i.XMLURL)
	}

	feeds, err := new(mongoModel.Feed).FindByUrls(c, feedUrls...)
	if err != nil {
		return
	}
	var needToSub []string
	var notNewFeedMap = make(map[string]struct{})
	for _, v := range feeds {
		needToSub = append(needToSub, v.ID)
		notNewFeedMap[v.Url] = struct{}{}
	}

	var user = mongoModel.User{
		ID: userID,
	}
	err = user.SubscribeFeeds(c, needToSub...)
	if err != nil {
		returnErrorMsg(c, errorhelper.WarpErr(infra.DataBaseError, err))
		return
	}

	go func() {
		ctx := context.TODO()
		needToSub = make([]string, 0)
		for _, i := range allItem {
			if _, has := notNewFeedMap[i.XMLURL]; has {
				continue
			}
			log.Info("feed", log.String("title", i.Title), log.String("url", i.XMLURL))
			feed, err := logic.Feed.GetFeed(ctx, i.XMLURL)
			if err != nil {
				log.Error("GetFeed", log.Err(err))
				continue
			}
			needToSub = append(needToSub, feed.ID)
		}
		err = user.SubscribeFeeds(ctx, needToSub...)
		if err != nil {
			log.Error("SubscribeFeeds", log.Err(err))
			return
		}
	}()
	returnSuccessMsg(c, "OK", nil)
}

func ExportOPML(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user = mongoModel.User{
		ID: userID.(string),
	}
	err := user.FindFeedByID(c)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	feeds, err := new(mongoModel.Feed).FindByIDs(c, user.SubFeeds)
	if err != nil {
		returnErrorMsg(c, infra.DataBaseError)
		return
	}
	userOPML := opml.OPML{
		Version: "2.0",
		Head: opml.Head{
			Title:       "OPML Export From Rssreader",
			DateCreated: time.Now().Format(time.RFC1123Z),
			Docs:        "http://www.opml.org/spec2",
		},
		Body: opml.Body{
			Outlines: []opml.Outline{},
		},
	}
	for _, i := range feeds {
		userOPML.Body.Outlines = append(userOPML.Body.Outlines, opml.Outline{Type: "rss", Text: i.Title, XMLURL: i.Url})
	}
	userXml, err := xml.Marshal(userOPML)
	if err != nil {
		returnErrorMsg(c, infra.ExportError)
		return
	}
	returnSuccessMsg(c, "OK", string(userXml))
}

func GetOutlinesItem(outline []opml.Outline) []*opml.Outline {
	var res []*opml.Outline
	for i := range outline {
		if len(outline[i].Outlines) != 0 {
			res = append(res, GetOutlinesItem(outline[i].Outlines)...)
		}
		res = append(res, &outline[i])
	}
	return res
}
