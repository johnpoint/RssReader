package apis

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"rssreader/src/model"
	"strings"
	"testing"
)

var config = model.Config{
	Debug:    false,
	TLS:      false,
	Database: "testDB.db",
	Salt:     "test",
	Port:     "1323",
}

func TestMain(m *testing.M) {
	log.Println("==== 初始化资源 ===")
	file, _ := os.Create("config.json")
	defer file.Close()
	databy, _ := json.Marshal(config)
	_, _ = io.WriteString(file, string(databy)) // 写入测试配置文件
	del := os.Remove("./testDB.db")
	if del != nil {
		fmt.Println("数据库已经初始化")
	}
	result := m.Run() //运行go的测试
	log.Println("=== 释放资源 ===")
	_ = os.Remove("./testDB.db")
	os.Exit(result) //退出程序
}

func TestAccessible(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "serverStatus", args: args{c: c}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Accessible(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Accessible() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckAuth(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name string
		args args
		want *model.JwtCustomClaims
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckAuth(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeedAsRead(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FeedAsRead(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("FeedAsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeedTodo(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FeedTodo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("FeedTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFeedList(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetFeedList(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetFeedList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetPostContent(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetPostContent(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetPostContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetPostList(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetPostList(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetPostList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetReadPostList(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetReadPostList(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetReadPostList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	u := model.User{}
	u.Mail = "test@test.test"
	u.Password = "testpassword"
	data := []byte(u.Mail + config.Salt + u.Password)
	has := md5.Sum(data)
	md5Password := fmt.Sprintf("%x", has)
	u.Password = md5Password
	_ = u.New()

	type args struct {
		mail     string
		password string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "passLogin", args: args{mail: "test@test.test", password: "testpassword"}, want: 200},
		{name: "unpassLogin", args: args{mail: "test@test.test", password: "wrongpassword"}, want: 0},
		{name: "notExistLogin", args: args{mail: "test2@test.test", password: "wrongpassword"}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			userJSON := `{"mail":"` + tt.args.mail + `","password":"` + tt.args.password + `"}`
			req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := Login(c)
			if err != nil {
				panic(err)
			}
			b := rec.Body.String()
			sj, err := simplejson.NewJson([]byte(b))
			if err != nil {
				panic(err)
			}
			if (rec.Body != nil) && sj.Get("code").MustInt64() != tt.want {
				t.Errorf("Login() got = %v, want %v", sj.Get("code").MustInt64(), tt.want)
			}
		})
	}
}

func TestPostAsRead(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostAsRead(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("PostAsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostAsUnRead(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostAsUnRead(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("PostAsUnRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Register(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchFeed(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SearchFeed(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SearchFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscribeFeed(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SubscribeFeed(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnSubscribeFeed(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnSubscribeFeed(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("UnSubscribeFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
