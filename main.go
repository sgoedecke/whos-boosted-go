package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Print("Could not read from $PORT. Setting port to 3000...")
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	log.Print("Server listening on port " + port)

	router.GET("/scan/:id", scanAccount)

	router.Run(":" + port)
}

func scanAccount(c *gin.Context) {
	id := c.Param("id")
	boostCheckData, err := openDotaLookup(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "could not fetch winrates",
		})
	}

	chance, reasons := boostChance(boostCheckData)

	c.JSON(200, gin.H{
		"chance":  chance,
		"reasons": reasons,
	})
}
