package depend

import (
	"RssReader/app/controller"
	"RssReader/config"
	jwtModel "RssReader/model/jwt"
	"RssReader/pkg/bootstrap"
	"context"
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Api struct{}

var _ bootstrap.Component = (*Api)(nil)

var identityKey = "user_id"

func (r *Api) Init(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	routerGin := gin.New()
	routerGin.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config.CORS,
		AllowMethods:     []string{"PUT", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "content-type", "Cookie", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routerGin.GET("/ping", controller.Pong)

	routerGin.GET("/api/syspost", controller.Pong) //服务器公告

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*jwtModel.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &jwtModel.User{
				UserID: claims[identityKey].(string),
			}
		},
		Authenticator: controller.Login,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			userData := data.(*jwtModel.User)
			if len(userData.UserID) != 0 {
				c.Set("user_id", userData.UserID)
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.JSON(http.StatusOK, controller.ApiResp{
				Code:    200,
				Message: message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		return err
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		return err
	}

	routerGin.POST("/api/login", authMiddleware.LoginHandler) //登录
	routerGin.POST("/api/register", controller.Register)      //注册

	w := routerGin.Group("/api")
	w.GET("/refresh_token", authMiddleware.RefreshHandler)
	w.Use(authMiddleware.MiddlewareFunc())
	f := w.Group("/feed")
	{
		f.GET("/list", controller.FeedList)                    // 列出已经订阅feed
		f.POST("/search", controller.SearchFeed)               // 搜索feed
		f.POST("/subscribe/:id", controller.SubscribeFeed)     // 订阅feed
		f.POST("/unsubscribe/:id", controller.UnSubscribeFeed) // 退订feed
		f.POST("/read/:id", controller.Pong)                   // 将该feed标为已读
	}
	p := w.Group("/post")
	{
		p.POST("/read/:id", controller.PostAsRead)       // 将文章设置为已读
		p.POST("/unread/:id", controller.PostAsUnRead)   // 将文章设置为未读
		p.POST("/readafter/", controller.PostAsUnRead)   // 稍后再读 TODO
		p.GET("/content/:id", controller.GetPostContent) // 获取文章内容
		p.GET("/:num", controller.PostList)              // 获取文章列表
		//p.GET("/", controller.Pong)                      // 获取文章列表
		p.GET("/read", controller.ReadPostList) // 获取已读文章列表
	}

	u := w.Group("/user")
	{
		u.POST("/password", controller.ChangePassword) // 重设密码
		u.GET("/opml", controller.ExportOPML)          // 导出opml
		u.POST("/opml", controller.ImportOPML)         // 导入opml
	}

	go func() {
		fmt.Println("[init] HTTP Listen at " + config.Config.HttpServerListen)
		err := routerGin.Run(config.Config.HttpServerListen)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}
