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
	app.POST("users", controllers.DeleteUsers)
	app.POST("user", controllers.InsertUser)
	app.POST("user/update", controllers.UpdateUser)
	app.POST("user/role/update", controllers.UpdateUserRole)

	app.GET("useritems/:id", controllers.GetActiveUserItems)

	app.GET("tags/:id", controllers.GetAllTags)
	app.POST("tag", controllers.UpdateUserTag)

	app.GET("tags", controllers.GetTags)
	app.POST("addFilter", controllers.CreateTag)
	app.DELETE("deleteFilter", controllers.DeleteTag)

	app.GET("itemTags/:id", controllers.GetAllItemTags)

	app.GET("items", controllers.GetItems)
	app.GET("/items/:id", controllers.GetItemByIdWithTheActiveRental)
	app.GET("/itemsDetail/:id", controllers.GetItemByIdWithAllRentals)
	app.POST("uploadItemPhoto/:id", controllers.UploadItemImage)
	app.GET("itemImage/:id", controllers.GetItemPhoto)
	app.POST("addItem", controllers.InsertItem)
	app.PUT("item", controllers.UpdateItem)
	app.DELETE("item/:id", controllers.DeleteItem)

	app.POST("startRental", controllers.InsertRental)
	app.POST("endRental", controllers.EndRental)
	app.GET("rentalhistory/:id", controllers.GetHistory)

	app.POST("report", controllers.InsertReport)
	app.GET("reports", controllers.GetReports)

	app.GET("/hello", func(c *gin.Context) { // bitte nicht löschen, ist gut zum testen
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

}
