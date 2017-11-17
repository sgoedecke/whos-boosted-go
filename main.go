package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	router.GET("/scan/:id", func(c *gin.Context) {
		id := c.Param("id")
		winrates, err := openDotaLookup(id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "could not fetch winrates",
			})
		}

		c.JSON(200, gin.H{
			"report": winrates,
		})
	})

	router.Run(":" + port)
}

var regionCodes = map[string]string{
	"0":  "AUTOMATIC",
	"1":  "US WEST",
	"2":  "US EAST",
	"3":  "EUROPE",
	"5":  "SINGAPORE",
	"6":  "DUBAI",
	"7":  "AUSTRALIA",
	"8":  "STOCKHOLM",
	"9":  "AUSTRIA",
	"10": "BRAZIL",
	"11": "SOUTHAFRICA",
	"12": "PW TELECOM SHANGHAI",
	"13": "PW UNICOM",
	"14": "CHILE",
	"15": "PERU",
	"16": "INDIA",
	"17": "PW TELECOM GUANGDONG",
	"18": "PW TELECOM ZHEJIANG",
	"19": "JAPAN",
	"20": "PW TELECOM WUHAN",
	"25": "PW UNICOM TIANJIN",
}

type OpenDotaReport struct {
	Region map[string]*RegionReport
}

type RegionReport struct {
	Games int `json: games`
	Win   int `json: win`
}

func openDotaLookup(id string) (map[string]int, error) {
	resp, _ := http.Get("https://api.opendota.com/api/players/" + id + "/counts")
	defer resp.Body.Close()
	var f OpenDotaReport
	err := json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return nil, err
	}

	wr := make(map[string]int)
	for regionCode, report := range f.Region {
		// if we've got enough games in a region, put it on the list
		if report.Games > 20 {
			rName := regionCodes[regionCode]
			wr[rName] = (report.Win * 100) / report.Games // wr expressed as percentage
		}
	}
	return wr, nil
}
