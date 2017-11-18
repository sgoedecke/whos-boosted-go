package main

import (
	"math"
	"strconv"
)

// calculates the chance of an account being boosted from a map of winrates
// by server. returns a percentage chance and a list of reasons
func boostChance(data *BoostCheckData) (int, []string) {
	reasons := make([]string, 0)
	chance := 0

	// get region with highest winrate
	highestWr := 0
	var highestWrRegion string
	for region, wr := range data.Winrates {
		if wr > highestWr {
			highestWr = wr
			highestWrRegion = region
		}
	}

	wrOnMostPlayedServer := data.Winrates[data.MostPlayedServer]
	wrDiff := float64(highestWr - wrOnMostPlayedServer)
	chance = int(math.Min(wrDiff*3, 99.0))
	wrDifference := "Winrate on most played region (" + data.MostPlayedServer + ", " + strconv.Itoa(wrOnMostPlayedServer) + "%) was " + strconv.FormatFloat(wrDiff, 'f', 1, 64) + " less than the highest winrate (" + highestWrRegion + ", " + strconv.Itoa(highestWr) + "%)"
	reasons = append(reasons, wrDifference)

	return chance, reasons
}
