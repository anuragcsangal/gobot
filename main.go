package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Token string

func init() {
	Token = "MTE4NTA2ODA2NDIxNjE4Njg4MA.G6zlZC.s-ApxpeTN3wCJ2H4oLnJm0zQ6YaPSYfIK4ueU4"
}

type imageUrl struct {
	URL string `json:"url"`
}

func getUrl() string {
	url := "https://api.thecatapi.com/v1/images/search"
	// url := "https://meowfacts.herokuapp.com/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []imageUrl
	json.Unmarshal([]byte(body), &result)

	return result[0].URL
}

func main() {

	// New discord session using provided bot token
	ds, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err)
	}

	// Register the messageCreate func as a callback for MessageCreate event
	ds.AddHandler(messageCreate)

	ds.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a web socket connection to Discord and begin listening
	err = ds.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is now running ...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the discord session
	ds.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	url := getUrl()
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "show-cat" {
		s.ChannelMessageSend(m.ChannelID, url)
	}
}
