package routes

import (
	"BorrowBox/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {

	app.GET("user/:id", controllers.UserById)
	app.GET("users", controllers.GetUsers)
	app.GET("items", controllers.GetItems)
	app.DELETE("user/:id", controllers.DeleteUser)
	app.POST("user", controllers.InsertUser)
	app.PUT("/user/:id", controllers.UpdateUser)
	app.GET("getDocumentByID/:collection/:id", controllers.GetDocumentByIDROute)
	app.POST("startRental", controllers.InsertRental)
	app.GET("useritems/:id", controllers.GetActiveUserItems)
	app.POST("login", controllers.Login)
	app.GET("/items/:id", controllers.GetItemByIdWithTheActiveRental)
	app.GET("/itemsDetail/:id", controllers.GetItemByIdWithAllRentals)

	app.GET("/hello", func(c *gin.Context) { // bitte nicht löschen, ist gut zum testen
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

}
