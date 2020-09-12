package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"rssreader/src/apis"
	"rssreader/src/model"
	"rssreader/src/spider"
)

func Run() {
	go func() {
		spider.Spider()
	}()
	conf := model.Config{}
	conf.Load()
	e := echo.New()
	e.Debug = conf.Debug
	e.HideBanner = true

	//echo中间件配置
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: conf.AllowOrigins,
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
	})) //CORS配置
	jwtConfig := middleware.JWTConfig{
		Claims:     &model.JwtCustomClaims{},
		SigningKey: []byte(conf.Salt),
	} //JsonWebToken 配置

	e.GET("/", apis.Accessible)            //服务端运行验证
	e.POST("/api/login", apis.Login)       //登录
	e.POST("/api/register", apis.Register) //注册

	w := e.Group("/api")
	w.Use(middleware.JWTWithConfig(jwtConfig))
	w.GET("/", apis.Accessible)
	f := w.Group("/feed")
	f.GET("/list", apis.GetFeedList)                 //列出已经订阅feed
	f.POST("/search", apis.SearchFeed)               //搜索feed
	f.POST("/subscribe/:id", apis.SubscribeFeed)     //订阅feed
	f.POST("/unsubscribe/:id", apis.UnSubscribeFeed) //退订feed
	f.POST("/read/:id", apis.FeedAsRead)             //将该feed标为已读
	p := w.Group("/post")
	p.POST("/read/:id", apis.PostAsRead)       //将文章设置为已读
	p.POST("/unread/:id", apis.PostAsUnRead)   //将文章设置为未读
	p.GET("/content/:id", apis.GetPostContent) //获取文章内容
	p.GET("/", apis.GetPostList)               //获取文章列表
	p.GET("/read", apis.GetReadPostList)       //获取已读文章列表

	if conf.TLS {
		e.Logger.Fatal(e.StartTLS(":"+conf.Port, conf.CERTPath, conf.KEYPath))
	} else {
		e.Logger.Fatal(e.Start(":" + conf.Port))
	}
}
