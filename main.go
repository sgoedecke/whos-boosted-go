package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sort"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Print("Could not read from $PORT. Setting port to 3000...")
		port = "3000"
	}

	router := gin.New()
	router.LoadHTMLGlob("templates/*")

	router.Use(gin.Logger())
	log.Print("Server listening on port " + port)

	router.GET("/", renderRoot)
	router.GET("/scan_friends", scanFriends)
	router.GET("/api/scan/:id", scanAccount)
	router.Run(":" + port)
}

func renderRoot(c *gin.Context) {
	c.HTML(200, "index.tmpl", nil)
}

type ScanResult struct {
	Id      string
	Name    string
	Chance  int
	Reasons []string
}

func scanFriends(c *gin.Context) {
	id := c.Query("id")
	friend_ids := getFriendIds(id)
	friend_hash := getNamesFromIds(friend_ids)

	resultsChan := make(chan ScanResult)
	var results []ScanResult

	for id, name := range friend_hash {
		go scanPlayer(id, name, resultsChan)
	}

	for i := 0; i < len(friend_hash); i++ {
		results = append(results, <-resultsChan)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Chance > results[j].Chance
	})

	c.HTML(200, "results.tmpl", gin.H{
		"results": results,
	})
}

func scanPlayer(id string, name string, results chan<- ScanResult) {
	boostCheckData, _ := openDotaLookup(id)
	chance, reasons := boostChance(boostCheckData)
	sr := ScanResult{id, name, chance, reasons}
	results <- sr
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
