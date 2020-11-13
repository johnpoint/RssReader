package apis

import (
	"encoding/json"
	"net/http"
	"rssreader/src/model"

	"github.com/labstack/echo/v4"
)

func GetFeedList(c echo.Context) error {
	u, err := CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Code: 0, Message: err.Error()})
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
		data = append(data, respFeed{ID: i.FID, Title: f.Title, Url: f.Url, Status: f.Status})
	}
	bdata, _ := json.Marshal(data)
	respdata := ""
	if string(bdata) == "null" {
		respdata = "[]"
	} else {
		respdata = string(bdata)
	}
	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: respdata})
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
