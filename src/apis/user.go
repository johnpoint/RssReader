package apis

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gilliek/go-opml/opml"
	"io"
	"log"
	"net/http"
	"rssreader/src/model"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type respFeed struct {
	ID         int64
	Title      string
	Url        string
	Subscriber int64
	Status     int64
}
type respPostList struct {
	ID        int64
	Feed      int64
	FeedTitle string
	Title     string
	Link      string
	Time      string
}

func ResetPassword(c echo.Context) error {
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	conf := model.Config{}
	err = conf.Load()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	salt := conf.Salt
	type resetPassword struct {
		Password string
	}
	p := resetPassword{}
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	data := []byte(u.Mail + salt + p.Password)
	has := md5.Sum(data)
	md5Password := fmt.Sprintf("%x", has)
	u.Password = md5Password
	err = u.Save()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func PostAsRead(c echo.Context) error {
	pid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	err = u.Read(pid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func FeedAsRead(c echo.Context) error {
	fid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	f := model.Feed{}
	f.ID = fid
	err = f.Get([]string{"id"})
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	p := f.Post(-1)
	for _, i := range p {
		err := u.Read(i.ID)
		if err != nil {
			return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
		}
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func PostAsUnRead(c echo.Context) error {
	pid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	err = u.UnRead(pid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func SubscribeFeed(c echo.Context) error {
	fid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	err = u.Subscribe(fid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func UnSubscribeFeed(c echo.Context) error {
	fid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	f := model.Feed{}
	f.ID = fid
	p := f.Post(-1)
	for _, i := range p {
		err := u.UnRead(i.ID)
		if err != nil {
			return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
		}
	}
	err = u.Unsubscribe(fid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func GetPostList(c echo.Context) error {
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	err, sub := u.Subscribes()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	getPostNum := c.Param("num")
	getPostNumI, err := strconv.ParseInt(getPostNum, 10, 64)
	if err != nil {
		getPostNumI = 50
	}
	if getPostNumI > 500 {
		getPostNumI = 50
	}
	var rep []respPostList
	var feedID []int64
	for _, i := range sub {
		feedID = append(feedID, i.FID)
	}
	p := model.Post{}
	items := p.FeedPost(feedID, int(getPostNumI))
	for _, i := range items {
		f := model.Feed{ID: i.FID}
		if err := f.Get([]string{"id", "title", "url"}); err != nil {
			return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
		}
		item := respPostList{
			ID:        i.ID,
			Feed:      i.FID,
			FeedTitle: f.Title,
			Link:      i.Url,
			Title:     i.Title,
			Time:      i.Published,
		}
		rep = append(rep, item)
	}
	sort.Slice(rep, func(i, j int) bool {
		if rep[i].Time > rep[j].Time {
			return true
		}
		return false
	})
	data, _ := json.Marshal(rep)
	respData := ""
	if string(data) == "null" {
		respData = "[]"
	} else {
		respData = string(data)
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: respData})
}

func GetReadPostList(c echo.Context) error {
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	list, _ := u.ReadPost()
	data, _ := json.Marshal(list)
	respData := ""
	if string(data) == "null" {
		respData = "[]"
	} else {
		respData = string(data)
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: respData})
}

func GetPostContent(c echo.Context) error {
	pid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	p := model.Post{ID: pid}
	err := p.Get(nil)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	b, err := json.Marshal(p)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(b)})
}

func ImportOPML(c echo.Context) error {
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	file, err := c.FormFile("opml")
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	defer src.Close()
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, src)
	// check errors
	err = u.Import(buf.String())
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func ExportOPML(c echo.Context) error {
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	err, sub := u.Subscribes()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
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
	for _, i := range sub {
		f := model.Feed{ID: i.FID}
		err := f.Get([]string{"id", "title", "url"})
		if err != nil {
			log.Println("feed not found")
			continue
		}
		userOPML.Body.Outlines = append(userOPML.Body.Outlines, opml.Outline{Type: "rss", Text: f.Title, XMLURL: f.Url})
	}
	userXml, err := xml.Marshal(userOPML)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Code: 0, Message: "Internal Server Error"})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(userXml)})
}
