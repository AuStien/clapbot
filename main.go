package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func main() {
	apiKey, ok := os.LookupEnv("CLAPBOT_SLACK_API_KEY")
	if !ok {
		log.Fatal("missing API key")
	}

	signingSecret, ok := os.LookupEnv("CLAPBOT_SLACK_SIGNING_SECRET")
	if !ok {
		log.Fatal("missing signing key")
	}

	api := slack.New(apiKey)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { return })
	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { return })

	http.HandleFunc("/clap", func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case "/clap":
			user, err := api.GetUserInfo(s.UserID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, _, err = api.PostMessage(s.ChannelID,
				slack.MsgOptionUsername(user.RealName),
				slack.MsgOptionIconURL(user.Profile.ImageOriginal),
				slack.MsgOptionText(addClap(s.Text), true))
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(200)
		case "/randomcase":
			user, err := api.GetUserInfo(s.UserID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, _, err = api.PostMessage(s.ChannelID,
				slack.MsgOptionUsername(user.RealName),
				slack.MsgOptionIconURL(user.Profile.ImageOriginal),
				slack.MsgOptionText(randomCase(s.Text, time.Now().Unix()), true))
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(200)
		default:
			fmt.Println("unknown command")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	fmt.Println("[INFO] Server listening on 8080")
	http.ListenAndServe(":8080", nil)
}

func addClap(text string) string {
	newText := ""
	for i, s := range text {
		// The clap never belongs at the start of the sentence
		if i == 0 {
			newText = fmt.Sprintf("%s%c", newText, s)
			continue
		}

		if i == len(text)-1 {
			newText = fmt.Sprintf("%s%c :clap:", newText, s)
			continue
		}

		if s == ' ' && (unicode.IsLetter(rune(text[i-1])) || unicode.IsPunct(rune(text[i-1]))) {
			newText = fmt.Sprintf("%s :clap:%c", newText, s)
		} else {
			newText = fmt.Sprintf("%s%c", newText, s)
		}
	}

	return newText
}

func randomCase(text string, seed int64) string {
	rand.Seed(seed)
	newText := ""
	for _, s := range text {
		randInt := rand.Intn(2)

		if randInt == 0 {
			newText = fmt.Sprintf("%s%s", newText, strings.ToUpper(string(s)))
		} else {
			newText = fmt.Sprintf("%s%s", newText, string(s))
		}
	}
	return newText
}
