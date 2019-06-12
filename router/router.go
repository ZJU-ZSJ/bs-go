package router

import (
	"bs-go/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {
	// Creates a default gin router
	r := gin.Default() // Grouping routes
	// group： api
	r.Use(cors.Default())
	api := r.Group("/api")
	{
		api.GET("/hello", handlers.HelloPage)
		api.POST("/register", handlers.RegisterPage)
		api.POST("/login", handlers.LoginPage)
	}
	book := r.Group("/book")
	{
		book.POST("/add", handlers.BookAdd)
		book.GET("/show/:id", handlers.Bookshow)
	}
	_ = r.Run(":8000") // listen and serve on 0.0.0.0:8000
}
