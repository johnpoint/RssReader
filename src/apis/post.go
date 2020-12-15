package apis

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"rssreader/src/model"
	"sort"
	"strconv"
)

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

func GetReadAfter(c echo.Context) error {
	type respRA struct {
		ID     int64
		Title  string
		Source string
	}
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	err, ra := u.GetReadAfter()
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
	}
	var r []respRA
	for _, i := range ra {
		p := model.Post{ID: i.PID}
		if err := p.Get([]string{"f_id", "title"}); err != nil {
			continue
		}
		f := model.Feed{ID: p.FID}
		if err := f.Get([]string{"title"}); err != nil {
			continue
		}
		r = append(r, respRA{ID: p.ID, Title: p.Title, Source: f.Title})
	}
	data, _ := json.Marshal(r)
	respData := ""
	if string(data) == "null" {
		respData = "[]"
	} else {
		respData = string(data)
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: respData})
}
