package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	slackVerificationToken := os.Getenv("SLACK_VERIFICATION_TOKEN")
	googleMapsApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	gurunaviKeyId := os.Getenv("GURUNAVI_KEY_ID")

	http.HandleFunc("/api/v1/meshi", func(w http.ResponseWriter, r *http.Request) {
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if !s.ValidateToken(slackVerificationToken) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		switch s.Command {
		case "/meshi":
			m := Meshi(
				googleMapsApiKey,
				35.6863929, 139.7004232,
				500,
				s.Text,
			)
			fmt.Println(m)

			params := m

			slack.NewPostMessageParameters()
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}
	})

	http.HandleFunc("/api/v1/meshi-detail", func(w http.ResponseWriter, r *http.Request) {
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if !s.ValidateToken(slackVerificationToken) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		switch s.Command {
		case "/meshi":
			m := MeshiDetail(
				gurunaviKeyId,
				"",
				"AREAM2115",
			)
			fmt.Println(m)

			params := m

			slack.NewPostMessageParameters()
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}
	})
	log.Println("[INFO] Server listening")
	http.ListenAndServe(":13071", nil)
}
