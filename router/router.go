package router

import (
	"bsgo/handlers"
	"github.com/gin-gonic/gin"
)

func Init() {
	// Creates a default gin router
	r := gin.Default() // Grouping routes
	r.Use(Cors())
	// group： api
	api := r.Group("/api")
	{
		api.GET("/hello", handlers.HelloPage)
		api.POST("/register", handlers.RegisterPage)
		api.POST("/login", handlers.LoginPage)
		api.GET("/info", handlers.Bookinfo)
		api.POST("/createorder", handlers.CreateOrder)
		api.GET("/myorder", handlers.MyOrder)
		api.GET("/mysell", handlers.MySell)
		api.POST("/bcom", handlers.Bcom)
		api.POST("/scom", handlers.Scom)
		api.GET("/userinfo", handlers.Userinfo)
		api.POST("/changeface", handlers.Changeface)
		api.GET("/search", handlers.Search)
	}
	book := r.Group("/book")
	{
		book.POST("/add", handlers.BookAdd)
		book.GET("/show/:id", handlers.Bookshow)
	}
	_ = r.Run(":8000") // listen and serve on 0.0.0.0:8000
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		// 处理js-ajax跨域问题
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Token")
		c.Next()
	}
}
