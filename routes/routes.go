package routes

import (
	"BorrowBox/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	app.POST("login", controllers.Login)

	app.GET("user/:id", controllers.UserById)
	app.GET("users", controllers.GetUsers)
	app.DELETE("user/:id", controllers.DeleteUser)
	app.POST("user", controllers.InsertUser)
	app.PUT("/user/:id", controllers.UpdateUser)

	app.GET("useritems/:id", controllers.GetActiveUserItems)

	app.GET("items", controllers.GetItems)
	app.GET("/items/:id", controllers.GetItemByIdWithTheActiveRental)
	app.GET("/itemsDetail/:id", controllers.GetItemByIdWithAllRentals)

	app.POST("startRental", controllers.InsertRental)
	app.PUT("endRental/:itemId", controllers.EndRental)

	app.GET("/hello", func(c *gin.Context) { // bitte nicht löschen, ist gut zum testen
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

}
