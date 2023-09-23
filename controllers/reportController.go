package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"fmt"
	Time "time"

	"github.com/gin-gonic/gin"
)

func InsertReport(c *gin.Context) {
	var report models.Report
	report.Time = Time.Now()
	c.BindJSON(&report)
	fmt.Println(report)
	_, err := database.InsertDocument("reports", report)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Report could not be inserted.",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Report inserted.",
		})
	}
}
