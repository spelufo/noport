package noport

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const configPath = "config.json"

func RunServer() {
	r := gin.Default()
	r.Use(middleware)
	r.GET("/config.json", handleConfigGet)
	r.POST("/config.json", handleConfigPost)
	r.Run()
}

func handleConfigGet(c *gin.Context) {
	b, err := os.ReadFile(configPath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Data(http.StatusOK, "application/json", b)
}

func handleConfigPost(c *gin.Context) {
	config := Config{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = os.WriteFile(configPath, file, 0644)
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


