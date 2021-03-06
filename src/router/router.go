package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"rssreader/src/apis"
	"rssreader/src/model"
	"rssreader/src/spider"
)

func Run(config string) {
	model.Init(config)
	go func() {
		spider.Spider()
	}()
	conf := model.Config{}
	err := conf.Load(config)
	if err != nil {
		panic("Fatal error! Configuration file failed to load")
	}
	e := echo.New()
	e.Debug = conf.Debug
	e.HideBanner = true

	//echo中间件配置
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: conf.AllowOrigins,
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
	})) //CORS配置
	jwtConfig := middleware.JWTConfig{
		Claims:       &model.JwtCustomClaims{},
		SigningKey:   []byte(conf.Salt),
		ErrorHandler: apis.JwtError,
	} //JsonWebToken 配置
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &model.SysContext{Context: c, Config: config}
			return next(cc)
		}
	}) //自定义 Context 配置

	e.GET("/", apis.Accessible)      //服务端运行验证
	e.GET("/api/syspost", apis.Post) //服务器公告

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
	p.POST("/readafter/", apis.GetReadAfter)   //将文章设置为未读
	p.GET("/content/:id", apis.GetPostContent) //获取文章内容
	p.GET("/:num", apis.GetPostList)           //获取文章列表
	p.GET("/", apis.GetPostList)               //获取文章列表
	p.GET("/read", apis.GetReadPostList)       //获取已读文章列表

	u := w.Group("/user")
	u.POST("/password", apis.ResetPassword) //重设密码
	u.GET("/opml", apis.ExportOPML)         //导出opml
	u.POST("/opml", apis.ImportOPML)        //导入opml

	if conf.TLS {
		e.Logger.Fatal(e.StartTLS(":"+conf.Port, conf.CERTPath, conf.KEYPath))
	} else {
		e.Logger.Fatal(e.Start(":" + conf.Port))
	}
}
