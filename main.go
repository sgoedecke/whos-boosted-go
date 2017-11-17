package main

import (
	"log"
  "net/http"
	"os"
	"github.com/gin-gonic/gin"
  "encoding/json"
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
    openDotaLookup(id)
		c.JSON(200, gin.H{
			"hello": id,
		})
	})

	router.Run(":" + port)
}

var regionCodes = map[string]string {
  "0": "AUTOMATIC",
  "1": "US WEST",
  "2": "US EAST",
  "3": "EUROPE",
  "5": "SINGAPORE",
  "6": "DUBAI",
  "7": "AUSTRALIA",
  "8": "STOCKHOLM",
  "9": "AUSTRIA",
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
  "25": "PW UNICOM TIANJIN"
}

type OpenDotaReport struct {
  Region map[string]*RegionReport
}

type RegionReport struct {
  Games int `json: games`
  Win int `json: win`
}

func openDotaLookup(id string) {
  resp, _ := http.Get("https://api.opendota.com/api/players/" + id + "/counts")
  defer resp.Body.Close()
  var f OpenDotaReport
  _ = json.NewDecoder(resp.Body).Decode(&f)

  var wr map[string]int
  for regionCode,report := range f.Region {
    // if we've got enough games in a region, put it on the list
    rName = regionCodes[regionCode]
    wr[rName] = (v.Win * 100) / v.Games // wr expressed as percentage
    log.Print(v.Games)
  }
}
