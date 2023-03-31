package noport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func runServer() {
	r := gin.Default()
	r.Use(middleware)
	r.GET("/"+theConfigFileName, handleConfigGet)
	r.POST("/"+theConfigFileName, handleConfigPost)
	r.Run()
}

func handleConfigGet(c *gin.Context) {
	c.JSON(http.StatusOK, theConfig)
}

func handleConfigPost(c *gin.Context) {
	config := Config{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := saveConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func middleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	if c.Request.Method == "OPTIONS" {
		c.Status(200)
		return
	}
}
