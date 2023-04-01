package noport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	theAddr = "localhost:8012"
	theConfigURL = "/.noport.json"
)

func runServer() {
	r := gin.Default()
	r.Use(middleware)
	r.GET(theConfigURL, handleConfigGet)
	r.POST(theConfigURL, handleConfigPost)
	r.POST("/install", handleInstallPost)
	r.Run(theAddr)
}

func handleConfigGet(c *gin.Context) {
	// Just load it each time, in case it is edited by hand.
	err := loadUserConfig()
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, theConfig)
}

func handleConfigPost(c *gin.Context) {
	config := Config{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	theConfig = config
	err := saveUserConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func handleInstallPost(c *gin.Context) {
	// Calling handleConfigPost is a bit of a hack. We want the install enpoint
	// to behave the same, saving the config first, and then installing it.
	// Assumes gin doesn't complain about setting the status more than once...
	handleConfigPost(c)

	err := installNginxConf(theNginxConfPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
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
