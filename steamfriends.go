package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const FRIENDS_ENDPOINT string = "http://api.steampowered.com/ISteamUser/GetFriendList/v1/"
const NAMES_ENDPOINT string = "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/"
const API_KEY string = "AD17383B71DE47BF52B2DD829548543C" // uh oh

// eg 76561198080304727

type FriendIdResponse struct {
	Friendslist struct {
		Friends []Friend `json:"friends"`
	}
}

type Friend struct {
	Steamid string `json:"steamid"`
}

func getFriendIds(id string) []string {
	req, _ := http.NewRequest("GET", FRIENDS_ENDPOINT, nil)
	q := req.URL.Query()
	q.Add("key", API_KEY)
	q.Add("steamid", id)
	q.Add("relationship", "friend")
	req.URL.RawQuery = q.Encode()
	log.Print(req.URL.String())
	resp, _ := http.Get(req.URL.String())
	defer resp.Body.Close()

	var f FriendIdResponse
	_ = json.NewDecoder(resp.Body).Decode(&f)

	ids := make([]string, 0)
	for _, friend := range f.Friendslist.Friends {
		ids = append(ids, friend.Steamid)
	}
	return ids
}

type NamesResponse struct {
	Response struct {
		Players []Player `json:"players"`
	}
}

type Player struct {
	Steamid string `json:"steamid"`
	Name    string `json:"personaname"`
}

// pass a hash of ids, get a hash of id/name back
func getNamesFromIds(ids []string) map[string]string {
	req, _ := http.NewRequest("GET", NAMES_ENDPOINT, nil)
	q := req.URL.Query()
	q.Add("key", API_KEY)
	q.Add("steamids", strings.Join(ids, ","))
	req.URL.RawQuery = q.Encode()
	log.Print(req.URL.String())
	resp, _ := http.Get(req.URL.String())
	defer resp.Body.Close()

	var nr NamesResponse
	_ = json.NewDecoder(resp.Body).Decode(&nr)

	result := make(map[string]string)

	for _, player := range nr.Response.Players {
		result[player.Steamid] = player.Name
	}

	return result
}
