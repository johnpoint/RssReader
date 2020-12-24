package apis

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"rssreader/src/model"
	"testing"
	"time"
)

var database = model.Database{
	Type:   "sqlite",
	DBname: "testDB.db",
}

var config = model.Config{
	Debug:    false,
	TLS:      false,
	Database: database,
	Salt:     "test",
	Port:     "1323",
}

func TestAccessible(t *testing.T) {
	tests := []struct {
		name    string
		wantErr int
	}{
		{name: "check_server", wantErr: 200},
	}
	for _, tt := range tests {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		t.Run(tt.name, func(t *testing.T) {
			if err := Accessible(c); rec.Code != tt.wantErr {
				fmt.Println(rec.Body)
				t.Errorf("Accessible() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckAuth(t *testing.T) {
	type args struct {
		user      model.User
		loginUser model.User
	}
	tests := []struct {
		name    string
		args    args
		want    model.User
		wantErr bool
	}{
		{name: "pass_auth", args: args{user: model.User{Mail: "t@t.t", Password: "123"}, loginUser: model.User{Mail: "t@t.t", Password: "123"}}, want: model.User{ID: 1, Mail: "t@t.t", Password: "123"}, wantErr: false},
		{name: "fail_auth", args: args{user: model.User{Mail: "t@t.t", Password: "123"}, loginUser: model.User{Mail: "t@t.t2", Password: "123"}}, want: model.User{}, wantErr: true},
	}
	for _, tt := range tests {
		file, _ := os.Create("config.json")
		defer file.Close()
		databy, _ := json.Marshal(config)
		_, _ = io.WriteString(file, string(databy))
		model.Init("config.json")
		u := tt.args.user
		_ = u.New()
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		claims := &model.JwtCustomClaims{
			Mail: tt.args.loginUser.Mail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 168).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckAuth(c)
			_ = os.Remove("testDB.db")
			_ = os.Remove("config.json")
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckAuth() got = %v, want %v", got, tt.want)
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
		file, _ := os.Create("config.json")
		defer file.Close()
		databy, _ := json.Marshal(config)
		_, _ = io.WriteString(file, string(databy))
		model.Init("config.json")
		t.Run(tt.name, func(t *testing.T) {
			if err := Register(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExportOPML(t *testing.T) {
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
			if err := ExportOPML(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExportOPML() error = %v, wantErr %v", err, tt.wantErr)
			}
			_ = os.Remove("testDB.db")
			_ = os.Remove("config.json")
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

func TestGetReadAfter(t *testing.T) {
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
			if err := GetReadAfter(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetReadAfter() error = %v, wantErr %v", err, tt.wantErr)
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

func TestImportOPML(t *testing.T) {
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
			if err := ImportOPML(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ImportOPML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJwtError(t *testing.T) {
	type args struct {
		in0 error
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
			if err := JwtError(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("JwtError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
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
			if err := Login(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPost(t *testing.T) {
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
			if err := Post(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
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

func TestResetPassword(t *testing.T) {
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
			if err := ResetPassword(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
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
