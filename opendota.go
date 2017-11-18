package main

import (
	"encoding/json"
	"net/http"
)

const MIN_REGION_GAMES int = 20

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

type BoostCheckData struct {
	Winrates         map[string]int
	MostPlayedServer string
}

// hits openDota API and returns a map of winrates by region
func openDotaLookup(id string) (*BoostCheckData, error) {
	resp, err := http.Get("https://api.opendota.com/api/players/" + id + "/counts")
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var f OpenDotaReport
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return nil, err
	}

	wr := make(map[string]int)
	var mostPlayed string
	mostGamesPlayed := 0
	for regionCode, report := range f.Region {
		// check if this region has the most games
		if report.Games > mostGamesPlayed {
			mostGamesPlayed = report.Games
			mostPlayed = regionCodes[regionCode]
		}

		// if we've got enough games in a region, put it on the list
		if report.Games > MIN_REGION_GAMES {
			rName := regionCodes[regionCode]
			wr[rName] = (report.Win * 100) / report.Games // wr expressed as percentage
		}
	}

	result := BoostCheckData{wr, mostPlayed}
	return &result, nil
}
