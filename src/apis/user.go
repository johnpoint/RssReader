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

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
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

func ResetPassword(c echo.Context) error {
	user := CheckAuth(c)
	conf := model.Config{}
	err := conf.Load()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	salt := conf.Salt
	u := model.User{Mail: user.Mail}
	err = u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
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

func CheckAuth(c echo.Context) *model.JwtCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	return claims
}

func PostAsRead(c echo.Context) error {
	pid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	err = u.Read(pid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func FeedAsRead(c echo.Context) error {
	fid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	f := model.Feed{}
	f.ID = fid
	err = f.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	p := f.Post()
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
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	err = u.UnRead(pid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

type respFeed struct {
	ID         int64
	Title      string
	Url        string
	Subscriber int64
}

func SearchFeed(c echo.Context) error {
	f := model.Feed{}
	if err := c.Bind(&f); err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	err := f.Get()
	if err != nil {
		err := f.New()
		if err != nil {
			return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
		}
		t, _ := json.Marshal(respFeed{ID: f.ID, Title: f.Title, Url: f.Url, Subscriber: f.Num})
		return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(t)})
	}
	t, _ := json.Marshal(respFeed{ID: f.ID, Title: f.Title, Url: f.Url, Subscriber: f.Num})
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(t)})
}

func SubscribeFeed(c echo.Context) error {
	fid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	err = u.AddSub(fid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

func UnSubscribeFeed(c echo.Context) error {
	fid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	f := model.Feed{}
	f.ID = fid
	p := f.Post()
	for _, i := range p {
		err := u.UnRead(i.ID)
		if err != nil {
			return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
		}
	}
	err = u.DelSub(fid)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: "OK"})
}

type respPostList struct {
	ID        int64
	Feed      int64
	FeedTitle string
	Title     string
	Link      string
	Time      string
}

func GetPostList(c echo.Context) error {
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	sub := u.Sub()
	getPostNum := c.Param("num")
	var rep []respPostList
	for _, i := range sub {
		f := model.Feed{ID: i.FID}
		f.Get()
		post := f.Post()
		for _, j := range post {
			item := respPostList{
				ID:        j.ID,
				Feed:      j.FID,
				FeedTitle: f.Title,
				Link:      j.Url,
				Title:     j.Title,
				Time:      j.Published,
			}
			rep = append(rep, item)
		}
	}
	sort.Slice(rep, func(i, j int) bool {
		if rep[i].Time > rep[j].Time {
			return true
		}
		return false
	})
	getPostNumI, err := strconv.ParseInt(getPostNum, 10, 64)
	if err != nil {
		getPostNumI = 50
	}
	if getPostNumI > 500 {
		getPostNumI = 50
	}
	if int64(len(rep)) >= getPostNumI {
		rep = rep[:getPostNumI]
	}
	data, _ := json.Marshal(rep)
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(data)})
}

func GetReadPostList(c echo.Context) error {
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	list, _ := u.ReadPost()
	data, _ := json.Marshal(list)
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(data)})
}

func GetPostContent(c echo.Context) error {
	pid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	p := model.Post{ID: pid}
	err := p.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	b, err := json.Marshal(p)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(b)})
}

func GetFeedList(c echo.Context) error {
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	sub := u.Sub()
	var data []respFeed
	for _, i := range sub {
		f := model.Feed{}
		f.ID = i.FID
		err := f.Get()
		if err != nil {
			return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
		}
		data = append(data, respFeed{ID: i.FID, Title: f.Title, Url: f.Url})
	}
	bdata, _ := json.Marshal(data)
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(bdata)})
}

func ImportOPML(c echo.Context) error {
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
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
	user := CheckAuth(c)
	u := model.User{Mail: user.Mail}
	err := u.Get()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "User does not exist"})
	}
	sub := u.Sub()
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
		err := f.Get()
		if err != nil {
			log.Println("feed not found")
			continue
		}
		userOPML.Body.Outlines = append(userOPML.Body.Outlines, opml.Outline{Type: "rss", Text: f.Title, XMLURL: f.Url})
	}
	userxml, err := xml.Marshal(userOPML)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Code: 0, Message: "Internal Server Error"})
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(userxml)})
}
