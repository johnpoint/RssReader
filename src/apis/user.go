package apis

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"github.com/gilliek/go-opml/opml"
	"io"
	"log"
	"net/http"
	"rssreader/src/model"
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
	cc := c.(*model.SysContext)
	err = conf.Load(cc.Config)
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
