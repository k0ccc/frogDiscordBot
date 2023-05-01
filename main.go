package main

import (
	"encoding/json"
	"fmt"

	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type redditUrl struct {
	Url string `json:"url"`
}

type redditPost []struct {
	Data redditUrl `json:"data"`
}

type redditData struct {
	Children redditPost `json:"children"`
}

type reddit struct {
	Data redditData `json:"data"`
}

var command string = "!wednsday"

func main() {
	dg, err := discordgo.New("Bot " + "")

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// getFrogs()
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == command+" ping" {
		s.ChannelMessageSend(m.ChannelID, " Pong!")
	}

	if m.Content == command+" pong" {
		s.ChannelMessageSend(m.ChannelID, " Ping!")
	}
	if m.Content == command {
		s.ChannelMessageSend(m.ChannelID, "what u want duude? if need somfing u may help this comnd: !wednsday help")
	}
	if m.Content == command+" test" {
		s.ChannelMessageSend(m.ChannelID, getFrogs())
	}
}
func randomIntForFrogs() int {
	randomInt := rand.Intn(24) + 1
	return randomInt
}

//	func randomIntForFrogs2() string {
//		r := rand.New(rand.NewSource(time.Now().UnixNano()))
//		return strconv.Itoa(r.Intn(25))
//	}
func getFrogs() string {
	// req
	site := "https://www.reddit.com/r/frog/hot.json"
	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	res, err := new(http.Client).Do(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer res.Body.Close()
	// reddit json parse
	var data reddit
	err = json.Unmarshal([]byte(b), &data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
	}
	frog := data.Data.Children[randomIntForFrogs()].Data.Url
	out := 0
	for {
		out += 1
		frogCut := frog[len(frog)-5:]
		// check if we get .jpeg, .jpg or something with 4 or 3 after dot
		// maybe we need check format of picture? improve?
		if frogCut[:1] == "." || frogCut[1:2] == "." {
			return frog
		} else {
			frog = data.Data.Children[randomIntForFrogs()].Data.Url
		}
		// BUG: because if randomIntForFrogs() gives us more then 2 repeating number, we will not reach the last post.
		if out > 24 {
			return ":frog:"
		}
	}
}
