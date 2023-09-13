package routes

import (
	"BorrowBox/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {

	app.GET("user/:id", controllers.UserById)
	app.GET("/:collection", controllers.GetDocuments)
	app.DELETE("user/:id", controllers.DeleteUser)
	app.POST("user", controllers.InsertUser)
	app.PUT("/user/:id", controllers.UpdateUser)
	app.GET("getDocumentByID/:collection/:id", controllers.GetDocumentByIDROute)
	app.POST("startRental", controllers.InsertRental)
	app.GET("useritems/:id", controllers.GetUserItems)
	app.POST("login", controllers.Login)
	app.GET("test/items/:id", controllers.GetItemByIdWithTheActiveRental)
	app.GET("testDetail/items/:id", controllers.GetItemByIdWithAllRentals)

	app.GET("/hello", func(c *gin.Context) { // bitte nicht löschen, ist gut zum testen
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

}
