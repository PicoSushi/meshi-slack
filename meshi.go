package main

import (
	"fmt"
	// "image"
	"log"
	"math/rand"
	"time"

	"github.com/nlopes/slack"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

const (
	FREE_QUERY_PER_SECOND = 2
)

var (
	PriceLevelMap = map[int]string{
		0: "安い",
		1: "普通",
		2: "高い",
		3: "かなり高い",
		4: "超高い",
	}
)

// Meshi searches some restraunt with given name and returns
func Meshi(api_key string, lat float64, lng float64, rad uint, keyword string) *slack.Msg {
	if keyword == "" {
		log.Printf("Bad keyword: %s", keyword)
		return &slack.Msg{Text: "検索ワードを指定してください。例: `/meshi パフェ`"}
	}

	c, err := maps.NewClient(maps.WithAPIKey(api_key), maps.WithRateLimit(FREE_QUERY_PER_SECOND))
	if err != nil {
		log.Printf("Client error: %s", err)
		return &slack.Msg{Text: "認証に失敗しました。ごめん。"}
	}

	r := &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: lat,
			Lng: lng,
		},
		Radius:  rad,
		Keyword: keyword,
		Type:    maps.PlaceTypeRestaurant,
		OpenNow: false,
	}

	response, err := c.NearbySearch(context.Background(), r)
	if err != nil {
		log.Printf("Search error: %s", err)
		return &slack.Msg{Text: fmt.Sprintf("「%s」の検索結果が見つかりませんでした。", keyword)}
	}
	restraunt := response.Results[rand.Intn(len(response.Results))]

	rand.Seed(time.Now().UnixNano())
	fmt.Println(restraunt.Name)

	a := slack.Attachment{}
	a.Fallback = restraunt.Name
	a.Pretext = fmt.Sprintf("これどうかな？「%s」の検索結果だよ。", keyword)
	a.Title = restraunt.Name
	a.Color = "#008000"

	a.Fields = []slack.AttachmentField{
		slack.AttachmentField{
			Title: "価格",
			Value: PriceLevelMap[restraunt.PriceLevel],
			Short: true,
		},
		slack.AttachmentField{
			Title: "評価",
			Value: fmt.Sprintf("%.1fツ星", restraunt.Rating),
			Short: true,
		},
		slack.AttachmentField{
			Title: "場所",
			Value: restraunt.Vicinity,
			Short: false,
		},
	}

	a.Footer = "/meshi command by @Ryota Kayanuma"

	msg := slack.Msg{}
	msg.ResponseType = "ephemeral"
	msg.Attachments = []slack.Attachment{a}
	return &msg
}
