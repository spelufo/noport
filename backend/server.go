package noport

import (
	// "fmt"
	"embed"
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	theAddr      = "localhost:8012"
	theConfigURL = "/.noport.json"
)

var devmode = os.Getenv("NOPORT_DEV") != ""

//go:embed public
var publicFS embed.FS

func runServer() {
	r := gin.Default()
	if devmode {
		r.Use(static.Serve("/", static.LocalFile("./resources/public", false)))
		r.Use(static.Serve("/js", static.LocalFile("./target/public/js", false)))
		r.Use(static.Serve("/cljs-out", static.LocalFile("./target/public/cljs-out", false)))
	} else {
		r.Use(static.Serve("/", EmbedFolder(publicFS, "public")))
	}
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
