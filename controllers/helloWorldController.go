package controllers

import (
	"BorrowBox/models"
	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {

	demo := models.Demo{
		Message: "Test",
	}

	c.JSON(200, demo)
}
