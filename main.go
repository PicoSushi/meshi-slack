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
				35.690921, 139.700258,
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
	log.Println("[INFO] Server listening")
	http.ListenAndServe(":13071", nil)
}
