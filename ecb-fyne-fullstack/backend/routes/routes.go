package routes

import (
	"backend/database"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Sensor endpoints
		api.POST("/sensor", createSensorData)
		api.GET("/sensor", getAllSensorData)
		api.GET("/sensor/latest", getLatestSensor)

		// Product endpoints
		api.GET("/products", getProducts)
		api.POST("/products", createProduct)

		// Transaction endpoints
		api.GET("/transactions", getTransactions)
		api.POST("/transactions", createTransaction)
	}
}

func createSensorData(c *gin.Context) {
	var data models.SensorData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data.Timestamp = time.Now()
	database.DB.Create(&data)
	c.JSON(http.StatusOK, data)
}

func getAllSensorData(c *gin.Context) {
	var data []models.SensorData
	database.DB.Order("timestamp desc").Limit(100).Find(&data)
	c.JSON(http.StatusOK, data)
}

func getLatestSensor(c *gin.Context) {
	var data models.SensorData
	result := database.DB.Order("timestamp desc").First(&data)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func getProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

func createProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&product)
	c.JSON(http.StatusOK, product)
}

func getTransactions(c *gin.Context) {
	var transactions []models.Transaction
	database.DB.Order("created_at desc").Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}

func createTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction.CreatedAt = time.Now()
	database.DB.Create(&transaction)
	c.JSON(http.StatusOK, transaction)
}
